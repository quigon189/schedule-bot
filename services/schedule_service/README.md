# Schedule Service

Микросервис для управления учебными расписаниями и изменениями в расписании. Предоставляет REST API для работы с расписаниями учебных групп и изменениями в расписании.

## Функционал

### Управление расписаниями групп
- Добавление, получение и удаление расписаний
- Фильтрация по учебному периоду, семестру и группе
- Поддержка различных форматов учебных лет (2025/2026)
- Хранение расписаний в виде изображений по URL

### Управление изменениями в расписании
- Добавление изменений на конкретные даты
- Получение изменений по дате (по умолчанию - текущая дата)
- Удаление изменений
- Поддержка множественных изображений для одного изменения

## API Endpoints

### Расписания групп
- `GET /api/v1/group_schedules` - получить расписание группы (фильтрация через query-параметры)
- `POST /api/v1/group_schedules` - добавить расписание для группы
- `DELETE /api/v1/group_schedules` - удалить расписание группы

### Изменения в расписании
- `GET /api/v1/changes` - получить изменения в расписании
- `POST /api/v1/changes` - добавить изменения в расписании
- `DELETE /api/v1/changes/{id}` - удалить изменения в расписании по ID

### Health check
- `GET /health` - проверка работоспособности сервиса

## Конфигурация

### Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `DB_HOST` | Хост базы данных | `localhost` |
| `DB_PORT` | Порт базы данных | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `""` |
| `DB_NAME` | Имя базы данных | `schedule_db` |

Сервер всегда запускается на порту `8080`.

### Пример .env файла
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secure_password
DB_NAME=schedule_db
```

## Запуск

### Требования
- Go 1.25+
- PostgreSQL 14+
- Docker (опционально)

### Docker
```bash
docker build -t schedule-service .
docker run -p 8080:8080 --env-file .env schedule-service
```

### Docker Compose
```yaml
# В составе общего проекта
docker-compose up schedule_service
```

## Миграции базы данных

Миграции находятся в директории `migrations/`:

1. `001_create_and_init_schedule_tables.sql` - создание таблиц расписаний
2. `002_add_group_schedules_view.sql` - создание представления для расписаний групп
3. `003_insert_test_data.sql` - тестовые данные
4. `004_group_schedules_2025-2026_1.sql` - расписания на 2025-2026 учебный год

При запуске сервиса миграции применяются автоматически через `goose`.

## Использование API

### Примеры запросов

#### 1. Получение расписания группы
```bash
curl -X GET "http://localhost:8080/api/v1/group_schedules?academic_year=2025/2026&half_year=1&group_name=СА-501" \
  -H "Content-Type: application/json"
```

#### 2. Добавление расписания группы
```bash
curl -X POST http://localhost:8080/api/v1/group_schedules \
  -H "Content-Type: application/json" \
  -d '{
    "academic_year": "2025/2026",
    "half_year": 1,
    "group_name": "СА-501",
    "semester": 5,
    "schedule_img_url": "https://example.com/schedule.jpg"
  }'
```

#### 3. Получение изменений в расписании
```bash
# Получение изменений на текущую дату
curl -X GET http://localhost:8080/api/v1/changes \
  -H "Content-Type: application/json"

# Получение изменений на конкретную дату
curl -X GET "http://localhost:8080/api/v1/changes?date=2025-09-01" \
  -H "Content-Type: application/json"
```

#### 4. Добавление изменений в расписании
```bash
curl -X POST http://localhost:8080/api/v1/changes \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2025-09-01",
    "image_urls": [
      "http://example.com/change1.jpg",
      "http://example.com/change2.jpg"
    ],
    "description": "Перенос лекции по математике"
  }'
```

#### 5. Удаление изменений
```bash
curl -X DELETE http://localhost:8080/api/v1/changes/1 \
  -H "Content-Type: application/json"
```

## Структура проекта

```
schedule_service/
├── cmd/server/                 # Точка входа
│   └── main.go                # Основной файл
├── internal/                   # Внутренние пакеты
│   ├── config/                # Конфигурация
│   ├── dto/                   # Data Transfer Objects (с валидацией)
│   ├── handlers/              # HTTP обработчики
│   ├── models/                # Модели данных
│   ├── repository/            # Работа с БД (репозитории)
│   └── service/               # Бизнес-логика
├── migrations/                # Миграции БД
├── pkg/utils/                 # Утилиты (ответы API)
└── docs/                      # Swagger документация
```

## Модели данных

### GroupSchedule
```go
type GroupSchedule struct {
    StudyPeriod    StudyPeriod
    Group          Group
    Semester       int
    ScheduleImgURL string
    CreatedAt      time.Time
}
```

### ScheduleChange
```go
type ScheduleChange struct {
    ID          int
    Date        time.Time
    ImgURLs     []string
    Description string
    CreatedAt   time.Time
}
```

## Валидация данных

Сервис использует библиотеку `validator/v10` для валидации входящих данных:

### Правила валидации:
- **Академический год:** формат `YYYY/YYYY` (регулярное выражение)
- **Название группы:** формат `БУКВЫ-ЦИФРЫ` (кириллица, регулярное выражение)
- **Дата:** формат `YYYY-MM-DD` (ISO)
- **URL:** стандартная валидация URL
- **Полугодие:** значение от 1 до 2
- **Семестр:** значение от 1 до 10

## База данных

### Основные таблицы:
1. `groups` - учебные группы
2. `study_periods` - учебные периоды (академический год, полугодие)
3. `group_schedules` - расписания групп
4. `changes` - изменения в расписании

### Представления:
- `group_schedules_view` - объединенное представление для удобного получения расписаний

## Особенности реализации

### Фильтрация расписаний
Сервис поддерживает гибкую фильтрацию через query-параметры. Можно комбинировать:
- Академический год
- Полугодие
- Название группы
- Семестр

### Изменения в расписании
- Поддерживается несколько изображений для одного изменения
- Изменения привязаны к конкретной дате
- При запросе без даты возвращаются изменения на текущую дату

### Производительность
- Используются представления базы данных для сложных запросов
- Поддерживается транзакционность при операциях записи
- Валидация выполняется до обращения к базе данных
