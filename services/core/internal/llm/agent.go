package llm

import (
	"context"
	"core/internal/models"
	"core/internal/services"
	"core/pkg/fileparser"
	"core/pkg/jsonschema"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Options map[string]any

type Client interface {
	Generate(ctx context.Context, systemPrompt, userMessage string, opts Options) (string, error)
	Chat(ctx context.Context, messages []Message, opts Options) (*Message, error)
}

type Message struct {
	Role      string
	Content   string
	Images    [][]byte
	ToolCalls []ToolCall
}

type ToolCall struct {
	Name      string
	Arguments json.RawMessage
}

type ToolDescription struct {
	Type     string `json:"type"`
	Function struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Parameters  json.RawMessage `json:"parameters"`
	} `json:"function"`
}

type LLMAgent struct {
	client       Client
	systemPrompt string
	tools        []Tool
	toolsByName  map[string]Tool
}

func NewAgent(client Client, svc services.ScheduleService, systemPrompt string) *LLMAgent {
	tools := initTools(svc)
	toolsByName := make(map[string]Tool)
	for _, t := range tools {
		toolsByName[t.Name] = t
	}
	return &LLMAgent{
		client:       client,
		systemPrompt: systemPrompt,
		tools:        tools,
		toolsByName:  toolsByName,
	}
}

func (a *LLMAgent) ProcessMessage(ctx context.Context, userMessage string) (string, error) {
	// 1. Формируем системный промпт с описанием инструментов

	tools, err := getToolDescriptions(a.tools)
	if err != nil {
		return "", fmt.Errorf("generate tool descriptions: %w", err)
	}

	now := time.Now()
	var userInfo strings.Builder
	user, ok := ctx.Value("user").(models.User)
	if ok {
		fmt.Fprintf(&userInfo, "%s ", user.FullName)
		roles := map[string]string{
			"user":    "Пользователь",
			"student": "Студент",
			"teacher": "Преподаватель",
			"manager": "Управляющий",
			"admin":   "Администратор",
		}
		for _, role := range user.Roles {
			fmt.Fprintf(&userInfo, "%s", roles[role.Name])
		}
		userInfo.WriteString("\n")
	}

	userPrompt := fmt.Sprintf(
		`Текущая дата: %s
Информация о пользователе:
%s
Сообщение:
%s`,
		now.Format("2006-01-02"), userInfo.String(), userMessage)

	log.Printf("Системная строка: %s", a.systemPrompt)
	log.Printf("Запрос пользователя: %s", userPrompt)

	messages := []Message{
		Message{
			Role:    "system",
			Content: a.systemPrompt,
		},
		Message{
			Role:    "user",
			Content: userPrompt,
		},
	}

	for i := 0; i < 10; i++ {
		respMessage, err := a.client.Chat(ctx, messages, Options{
			"temperature": 0,
			"tools":       tools,
		})
		if err != nil {
			return "", fmt.Errorf("LLM get message failed: %w", err)
		}
		jsonTC, _ := json.MarshalIndent(respMessage, "", "  ")
		log.Printf("Сообщение от LLM:\n%s", jsonTC)

		messages = append(messages, *respMessage)

		if len(respMessage.ToolCalls) > 0 {
			for _, tc := range respMessage.ToolCalls {
				var content string
				tool, exists := a.toolsByName[tc.Name]
				if exists {
					content, err = tool.Handler(ctx, tc.Arguments)
					if err != nil {
						content = fmt.Sprintf("Ошибка при выполнении инструмента: %v", err)
					}
				} else {
					content = fmt.Sprintf("Инструмент %s не найден", tc.Name)
				}

				toolMessage := Message{
					Role:    "tool",
					Content: content,
				}
				jsonToolMes, _ := json.MarshalIndent(toolMessage, "", "  ")
				log.Printf("Tool message: %s", jsonToolMes)
				messages = append(messages, toolMessage)
			}
		} else {
			return respMessage.Content, nil
		}
	}

	return "Извените не удалось получить ответ", nil
}

func (a *LLMAgent) GenerateStructuredData(ctx context.Context, userMessage string, parsedContent []fileparser.ParsedContent, target any) error {
	targetSchema, err := jsonschema.GenerateSchema(target)
	if err != nil {
		return fmt.Errorf("generate jsonschema: %w", err)
	}

	systemPrompt := fmt.Sprintf(`Ты – помощник по заполнению шаблона расписания. Из предоставленных данных (текст или изображения) извлеки информацию о группе, дисциплинах и студентах.
Обязательно используй инструмент set_group_info!
Если каких-то данных нет, оставь пустой массив или пустую строку. Для дат используй формат YYYY-MM-DD.`, string(targetSchema))

	var userPrompt strings.Builder
	var images [][]byte
	userPrompt.WriteString("Сообщение от пользователя:\n")
	userPrompt.WriteString(userMessage + "\n")
	for i, pc := range parsedContent {
		userPrompt.WriteString(fmt.Sprintf("Файл %d:\n", i+1))
		userPrompt.WriteString(pc.Text)
		for _, img := range pc.Images {
			images = append(images, img.Data)
		}
	}
	userPrompt.WriteString("Обязательно вызови инструмент set_group_info!!!")

	var messages []Message
	messages = append(messages, Message{
		Role:    "system",
		Content: systemPrompt,
	})
	messages = append(messages, Message{
		Role:    "user",
		Content: userPrompt.String(),
		Images:  images,
	})

	log.Printf("Messages: %+v", messages)

	tool := Tool{
		Name:        "set_group_info",
		Description: "Обязательно вызови эту функцию, для отправки данных пользователю",
		Parameters:  target,
		Handler: func(ctx context.Context, params json.RawMessage) (string, error) {
			if err := json.Unmarshal(params, &target); err != nil {
				return "", err
			}
			log.Printf("%+v", target)
			return "ok", nil
		},
	}

	tools := []Tool{tool}
	td, err := getToolDescriptions(tools)
	if err != nil {
		return err
	}

	respMessage, err := a.client.Chat(ctx, messages, Options{
		"temperature": 0,
		"tools":       td,
		"num_ctx":     32768,
	})
	if err != nil {
		return fmt.Errorf("generate llm response: %w", err)
	}

	log.Printf("LLM answer: %s", respMessage.Content)

	if len(respMessage.ToolCalls) > 0 {
		for _, tc := range respMessage.ToolCalls {
			if tc.Name == "set_group_info" {
				_, err := tool.Handler(ctx, tc.Arguments)
				if err != nil {
					return err
				} else {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("tool call error")
}

func getToolDescriptions(tools []Tool) ([]ToolDescription, error) {
	var toolDescriptions []ToolDescription
	for _, tool := range tools {
		var td ToolDescription
		td.Type = "function"
		td.Function.Name = tool.Name
		td.Function.Description = tool.Description
		schema, err := jsonschema.GenerateSchema(tool.Parameters)
		if err != nil {
			return nil, err
		}
		td.Function.Parameters = schema

		toolDescriptions = append(toolDescriptions, td)
	}

	return toolDescriptions, nil
}
