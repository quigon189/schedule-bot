# MCP Server

Model Context Protocol (MCP) сервер для интеграции с системой расписания. Предоставляет инструменты для AI-ассистентов и других клиентов для взаимодействия с системой расписаний, изменений в расписании и отправки уведомлений в Telegram.

## Функционал

### Инструменты MCP:
1. **Отправка изменений пользователю** - отправка уведомлений об изменениях в Telegram

## Конфигурация

### Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `SCHEDULE_SERVICE_URL` | URL сервиса расписаний | `http://schedule-service:8080` |
| `TG_SERVICE_URL` | URL Telegram бот-сервиса | `http://telegram-bot:8000` |
| `TIMEOUT` | Таймаут HTTP запросов (в секундах) | `30` |
| `SERVER_PORT` | Порт MCP сервера | `8080` |

### Пример .env файла
```env
SCHEDULE_SERVICE_URL=http://schedule-service:8080
TG_SERVICE_URL=http://bot-service:8000
TIMEOUT=30
```

## Запуск

### Требования
- Go 1.25+
- Docker (опционально)
- Работающие сервисы Schedule Service и Telegram Bot Service

### Docker
```bash
docker build -t mcp-server .
docker run -p 8080:8080 --env-file .env mcp-server
```

### Docker Compose
```yaml
# В составе общего проекта
docker-compose up mcp_server
```

## Структура проекта

```
mcp_server/
├── cmd/server/                 # Точка входа
│   └── main.go                # Основной файл
├── internal/                   # Внутренние пакеты
│   ├── config/                # Конфигурация
│   ├── models/                # Модели данных
│   ├── server/                # MCP и HTTP сервер
│   ├── service/               # Сервисы для внешних API
│   └── tools/                 # Инструменты MCP и их обработчики
├── go.mod                     # Зависимости Go
├── go.sum                     # Контрольные суммы зависимостей
└── Dockerfile                 # Конфигурация Docker
```

## Инструменты MCP

### 1. `send_changes_to_user`
**Описание:** Отправляет пользователю изменения в расписании через Telegram бота

**Параметры:**
- `date` (string, required): Дата в формате ISO 8601 YYYY-MM-DD
- `chat_id` (string, required): chat_id пользователя в Telegram

**Пример запроса:**
```json
{
  "name": "send_changes_to_user",
  "arguments": {
    "date": "2025-09-01",
    "chat_id": "123456789"
  }
}
```

## Состояние разработки

### Активная разработка:
⚠️ **Внимание:** Сервер находится в активной разработке. Возможны изменения в:
- Структуре инструментов
- API внешних сервисов
- Моделях данных
