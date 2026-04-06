package llm

import (
	"context"
	"core/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

	systemPrompt := fmt.Sprintf(`Ты — ассистент по расписанию занятий. Твоя задача — понять, какую информацию хочет получить пользователь, и выбрать подходящий инструмент из списка ниже.

## Доступные инструменты:
%s

## Правила:
1. Если пользователь спрашивает о расписании (сегодня, завтра, на неделю), используй get_group_schedule
2. Если вопрос не связан с расписанием, ответь: {"action": "none", "params": {}}

## Важно:
- Извлекай параметры из вопроса пользователя
- Если параметр не указан, оставь его пустым или попроси отправить новое сообщение, но с конкретными уточнениями
- Для дат используй формат YYYY-MM-DD

Ответь ТОЛЬКО JSON-объектом в формате:
{"action": "имя_инструмента", "params": {"параметр": "значение"}}`, string(toolsJSON))

	log.Printf("Запрос пользователя: %s", userMessage)
	log.Printf("Системная строка: %s", systemPrompt)

	// 2. Получаем решение от LLM
	actionResp, err := a.client.Generate(ctx, systemPrompt, userMessage, Options{"temperature": 0.1})
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
		Action string                 `json:"action"`
		Params map[string]interface{} `json:"params"`
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
	return a.formatResponse(ctx, userMessage, result)
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
- Форматируй ответ для удобного чтения (можно использовать переносы строк)

Ответ:`, data)

	return a.client.Generate(ctx, systemPrompt, userMessage, Options{"temperature": 0.2})
}
