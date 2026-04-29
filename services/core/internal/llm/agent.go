package llm

import (
	"context"
	"core/internal/models"
	"core/internal/services"
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
	client      Client
	tools       []Tool
	toolsByName map[string]Tool
}

func NewAgent(client Client, svc services.ScheduleService) *LLMAgent {
	tools := initTools(svc)
	toolsByName := make(map[string]Tool)
	for _, t := range tools {
		toolsByName[t.Name] = t
	}
	return &LLMAgent{
		client:      client,
		tools:       tools,
		toolsByName: toolsByName,
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

	systemPrompt := fmt.Sprintf(`Ты — ассистент по расписанию занятий. Твоя задача — понять, какую информацию хочет получить пользователь, выбрать подходящий инструмент из списка и дать пользователю короткий ответ.
Отвечай только на вопросы связанные с расписанием заняти.`)

	userPrompt := fmt.Sprintf(
		`Текущая дата: %s
Информация о пользователе:
%s
Сообщение:
%s`,
		now.Format("2006-01-02"), userInfo.String(), userMessage)

	log.Printf("Системная строка: %s", systemPrompt)
	log.Printf("Запрос пользователя: %s", userPrompt)

	messages := []Message{
		Message{
			Role:    "system",
			Content: systemPrompt,
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

// directAnswer - прямой ответ от LLM без вызова инструментов
func (a *LLMAgent) directAnswer(ctx context.Context, userMessage string) (string, error) {
	systemPrompt := `Ты — дружелюбный ассистент. Ответь пользователю на его вопрос, но только если он связан с учёбой, расписанием или учебным процессом. Если вопрос не по теме, вежливо скажи, что ты помогаешь только с расписанием.`
	return a.client.Generate(ctx, systemPrompt, userMessage, Options{"temperature": 0.2})
}

// formatResponse - форматирует результат работы инструмента в понятный ответ
func (a *LLMAgent) formatResponse(ctx context.Context, userMessage string, data string) (string, error) {
	systemPrompt := fmt.Sprintf(`Ты — ассистент по расписанию.
На запрос пользователя система предоставила следующие данные:
%s

Задача: сформулируй краткий, понятный и дружелюбный ответ для пользователя на русском языке.
- Не упоминай ID, JSON или технические детали
- Если данных много, выдели самое важное
- Если данных нет, скажи об этом вежливо
- Форматируй ответ для удобного чтения (можно использовать переносы строк)`, data)

	log.Printf("Системное сообщение для ответа пользователю:\n%s", systemPrompt)
	return a.client.Generate(ctx, systemPrompt, userMessage, Options{"temperature": 0.2})
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
