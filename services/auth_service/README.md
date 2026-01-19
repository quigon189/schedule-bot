# Auth Service

Микросервис для управления аутентификацией и авторизацией пользователей в системе Schedule Bot. Предоставляет REST API для управления пользователями, ролями и регистрационными кодами.

## Функционал

### Управление пользователями
- Создание, получение, обновление и удаление пользователей
- Управление ролями пользователей
- Регистрация через Telegram ID
- Привязка студентов к учебным группам

### Регистрационные коды
- Генерация уникальных кодов для регистрации
- Ограничение по времени действия и количеству использований
- Ролевое распределение (студент, преподаватель, менеджер)
- Привязка студентов к группам

### Роли
- **admin** - администратор системы
- **manager** - менеджер (может создавать коды)
- **teacher** - преподаватель
- **student** - студент (с привязкой к группе)
- **user** - базовый пользователь

## API Endpoints

### Пользователи
- `POST /api/v1/users` - создать пользователя
- `GET /api/v1/users/{telegram_id}` - получить пользователя
- `PUT /api/v1/users/{telegram_id}` - обновить пользователя
- `DELETE /api/v1/users/{telegram_id}` - удалить пользователя
- `POST /api/v1/users/register` - зарегистрироваться по коду

### Коды доступа
- `POST /api/v1/code/create` - создать регистрационный код

## Конфигурация

### Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `DB_HOST` | Хост базы данных | `localhost` |
| `DB_PORT` | Порт базы данных | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `postgres` |
| `DB_NAME` | Имя базы данных | `auth_service` |
| `SERVER_PORT` | Порт сервиса | `8080` |
| `CODE_CHARSET` | Символы для генерации кодов | `ABCDEFGHJKLMNPQRSTUVWXYZ23456789` |
| `CODE_LENGTH` | Длина кода | `6` |
| `MAX_GENERATE_TRIES` | Макс. попыток генерации уникального кода | `10` |
| `ADMINS` | Список Telegram ID администраторов (через пробел) | `1` |

### Пример .env файла
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=auth_service
SERVER_PORT=8080
ADMINS=123456789 987654321
```

## Запуск

### Требования
- Go 1.25+
- PostgreSQL 14+
- Docker (опционально)

### Docker
```bash
docker build -t auth-service .
docker run -p 8080:8080 --env-file .env auth-service
```

### Docker Compose
```yaml
# В составе общего проекта
docker-compose up auth_service
```

## Миграции базы данных

Миграции находятся в директории `migrations/`:

1. `001_create_user_tables.sql` - создание таблиц пользователей и ролей
2. `002_create_registration_codes.sql` - создание таблиц регистрационных кодов

При запуске сервиса миграции применяются автоматически через `goose`.

## Использование API

### Примеры запросов

#### 1. Создание пользователя
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "telegram_id": 123456789,
    "username": "john_doe",
    "full_name": "John Doe"
  }'
```

#### 2. Создание регистрационного кода
```bash
curl -X POST http://localhost:8080/api/v1/code/create \
  -H "Content-Type: application/json" \
  -d '{
    "role_name": "student",
    "group_name": "СА-501",
    "max_uses": 10,
    "expiration": 3600,
    "created_by": 1
  }'
```

#### 3. Регистрация по коду
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "code": "ABC123",
    "telegram_id": 123456789,
    "username": "jane_doe",
    "full_name": "Jane Doe"
  }'
```

## Структура проекта

```
auth_service/
├── cmd/server/                 # Точка входа
│   ├── main.go                # Основной файл
│   └── init_admins.go         # Инициализация админов
├── internal/                   # Внутренние пакеты
│   ├── config/                # Конфигурация
│   ├── dto/                   # Data Transfer Objects
│   ├── handlers/              # HTTP обработчики
│   ├── models/                # Модели данных
│   ├── repository/            # Работа с БД
│   └── service/               # Бизнес-логика
├── migrations/                # Миграции БД
├── pkg/utils/                 # Утилиты
└── docs/                      # Swagger документация
```

## Модели данных

### User
```go
type User struct {
    TelegramID int64
    Username   string
    FullName   string
    CreatedAt  time.Time
    UpdatedAt  time.Time
    IsActive   bool
    Roles      []Role
    Student    *Student
}
```

### RegistrationCode
```go
type RegistrationCode struct {
    ID          int
    Code        string
    RoleID      int
    GroupName   *string
    MaxUses     int
    CurrentUses int
    CreatedBy   *int64
    ExpiresAt   time.Time
    CreatedAt   time.Time
    Role        Role
    Creater     *User
}
```

## Swagger документация

После запуска сервиса документация доступна по адресу:
- `http://localhost:8080/swagger/index.html`
- `http://localhost:8080/swagger/doc.json`
