package llm

import (
	"context"
	"core/internal/models"
	"core/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Options map[string]any

type Client interface {
	Generate(ctx context.Context, systemPrompt, userMessage string, opts Options) (string, error)
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

	toolsJSON, _ := json.MarshalIndent(a.tools, "", "  ")

	now := time.Now()
	var userInfo strings.Builder
	user, ok := ctx.Value("user").(models.User)
	if ok {
		fmt.Fprintf(&userInfo, "Полное имя: %s\n", user.FullName)
		fmt.Fprintf(&userInfo, "Роли:")
		roles := map[string]string{
			"user":    "Пользователь",
			"student": "Студент",
			"teacher": "Преподаватель",
			"manager": "Управляющий",
			"admin":   "Администратор",
		}
		for _, role := range user.Roles {
			fmt.Fprintf(&userInfo, " %s", roles[role.Name])
		}
		userInfo.WriteString("\n")
	}

	systemPrompt := fmt.Sprintf(`Ты — ассистент по расписанию занятий. Твоя задача — понять, какую информацию хочет получить пользователь, и выбрать подходящий инструмент из списка ниже.

## Доступные инструменты:
%s

## Дополниетльная информация:
- Текущая дата: %s
- Подставляй фамилию преподавателя из информации о пользователе, только если тебя явно попросили (например в сообщении указано покажи мое расписание)

## Правила:
1. Если пользователь спрашивает о расписании (сегодня, завтра, на неделю), используй get_group_schedule
2. Если вопрос не связан с расписанием, ответь: {"action": "none", "params": {}}

## Важно:
- Извлекай параметры из вопроса пользователя
- Если параметр не указан, оставь его пустым
- Для дат используй формат YYYY-MM-DD

Ответь ТОЛЬКО JSON-объектом в формате:
{"action": "имя_инструмента", "params": {"параметр": "значение"}}`, string(toolsJSON), now.Format("2006-01-02"))

	userPrompt := fmt.Sprintf(`Информация о пользователе:
%s
Сообщение:
%s`, userInfo.String(), userMessage)

	log.Printf("Системная строка: %s", systemPrompt)
	log.Printf("Запрос пользователя: %s", userPrompt)

	// 2. Получаем решение от LLM
	actionResp, err := a.client.Generate(ctx, systemPrompt, userPrompt, Options{"temperature": 0})
	if err != nil {
		return "", fmt.Errorf("LLM action selection failed: %w", err)
	}

	log.Printf("Полученые действия от LLM: %v", actionResp)

	// Очищаем ответ
	actionResp = strings.TrimSpace(actionResp)
	actionResp = strings.TrimPrefix(actionResp, "```json")
	actionResp = strings.TrimSuffix(actionResp, "```")
	actionResp = strings.TrimSpace(actionResp)

	var decision struct {
		Action string         `json:"action"`
		Params map[string]any `json:"params"`
	}
	if err := json.Unmarshal([]byte(actionResp), &decision); err != nil {
		// Если не распарсили, пробуем ответить напрямую через LLM
		return a.directAnswer(ctx, userMessage)
	}

	// 3. Если инструмент не выбран, отвечаем напрямую
	if decision.Action == "none" {
		return a.directAnswer(ctx, userMessage)
	}

	// 4. Находим и выполняем инструмент
	tool, exists := a.toolsByName[decision.Action]
	if !exists {
		return a.directAnswer(ctx, userMessage)
	}

	// Выполняем инструмент
	result, err := tool.Handler(ctx, decision.Params)
	if err != nil {
		// Если ошибка, дадим понятный ответ
		return fmt.Sprintf("Извините, не удалось получить информацию: %v", err), nil
	}

	// 5. Формируем финальный ответ на основе полученных данных
	return a.formatResponse(ctx, userPrompt, result)
	// Тестово пробуем вернуть ответ пользователю
	//return result, nil
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
