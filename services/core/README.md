## Введение

Сервис предоставляет REST API для управления пользователями, ролями, группами, аудиториями, предметами, расписанием и журналом занятий. Аутентификация осуществляется с помощью JWT токенов.

**Базовый URL**: `http://<host>:<port>` (по умолчанию порт `8080`)

**Формат ответов**: JSON

**Структура ответа**:
- Успешный ответ:
  ```json
  {
    "success": true,
    "message": "строка",
    "data": { ... } // может отсутствовать
  }
  ```
- Ошибочный ответ:
  ```json
  {
    "success": false,
    "error": "строка ошибки"
  }
  ```

**Аутентификация**:  
Для доступа к защищённым эндпоинтам требуется передавать JWT access-токен в заголовке:  
`Authorization: Bearer <access_token>`

**Роли**:
- `admin` – полный доступ ко всем административным операциям.
- `manager` – управление ролями (назначение/удаление ролей `admin`, `manager`).
- `teacher` – преподаватель.
- `student` – студент с привязкой к группе.
- `user` – базовая роль, присваивается при регистрации.

**Коды ответов**:
- `200 OK` – успешное выполнение.
- `201 Created` – ресурс создан.
- `400 Bad Request` – ошибка валидации или неверный формат.
- `401 Unauthorized` – отсутствует или невалидный JWT токен.
- `403 Forbidden` – недостаточно прав.
- `404 Not Found` – ресурс не найден.
- `500 Internal Server Error` – внутренняя ошибка сервера.

---

## Эндпоинты аутентификации

### POST `/login`
Вход в систему. Получение access и refresh токенов, создание сессии.

**Тело запроса**:
```json
{
  "username": "string",   // обязательное, 3–20 символов (буквы, цифры, _)
  "password": "string",   // обязательное
  "user_agent": "string", // опционально (заполнится из заголовка User-Agent)
  "client_ip": "string"   // опционально (заполнится из X-Forwarded-For или RemoteAddr)
}
```

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "success login",
  "data": {
    "access_token": "string",
    "refresh_token": "string",
    "session_id": "string",
    "roles": ["admin", "user"]
  }
}
```

**Ошибки**:  
- `400` – ошибка валидации или неверные учётные данные.  
- `500` – внутренняя ошибка.

---

### POST `/refresh`
Обновление access-токена с использованием refresh-токена.

**Тело запроса**:
```json
{
  "session_id": "string",   // обязательное
  "refresh_token": "string" // обязательное
}
```

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "success refresh token",
  "data": {
    "access_token": "string",
    "refresh_token": "string"
  }
}
```

**Ошибки**:  
- `400` – неверный формат.  
- `500` – внутренняя ошибка.

---

### GET `/logout`
Завершение текущей сессии. Требует валидный access-токен.

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "success logout",
  "data": null
}
```

**Ошибки**:  
- `401` – невалидный токен.  
- `500` – внутренняя ошибка.

---

## Административные эндпоинты (требуется роль `admin`)

### GET `/admin/sessions`
Получить список всех активных сессий.

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "success get sessions",
  "data": [ ... ] // массив Session
}
```

**Ошибки**:  
- `403` – недостаточно прав.  
- `500` – внутренняя ошибка.

---

### Управление пользователями

#### GET `/admin/users`
Получить список пользователей с пагинацией.

**Параметры query**:
- `page` (int, опционально, по умолчанию 1) – номер страницы.
- `per_page` (int, опционально, по умолчанию 20) – записей на странице.
- `sort_by` (string, опционально, `id`, `username`, `full_name`, `email`, `created_at`).
- `sort_order` (string, опционально, `ASC` или `DESC`).

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "paginated users",
  "data": {
    "users": [ ... ],
    "total": 100,
    "page": 1,
    "per_page": 20,
    "total_pages": 5
  }
}
```

#### GET `/admin/users/{id}`
Получить информацию о пользователе по ID.

**Параметры пути**: `id` (int)

**Успешный ответ (200 OK)** – объект `User`.

#### POST `/admin/users`
Создать нового пользователя (базовая роль `user`).

**Тело запроса**:
```json
{
  "username": "string",  // 3–20, буквы/цифры/_
  "full_name": "string", // мин. 3 символа
  "email": "string",     // валидный email
  "password": "string"   // мин. 6 символов
}
```

**Успешный ответ (201 Created)** – объект `User`.

#### PATCH `/admin/users/password`
Сменить пароль пользователя (администратором).

**Тело запроса**:
```json
{
  "user_id": 123,        // обязательное
  "new_password": "string" // мин. 6 символов
}
```

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "password changed",
  "data": null
}
```

---

### Управление ролями

#### POST `/admin/roles/assign`
Назначить роль `admin` или `manager`.

**Тело запроса**:
```json
{
  "user_id": 123,
  "role_name": "admin"   // admin или manager
}
```

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "role assigned",
  "data": null
}
```

#### DELETE `/admin/roles/remove`
Удалить роль у пользователя.

**Тело запроса**:
```json
{
  "user_id": 123,
  "role_name": "student" // admin, manager, student, teacher
}
```

**Успешный ответ (200 OK)** – аналогично.

---

## Эндпоинты текущего пользователя (требуется аутентификация)

### GET `/user`
Получить информацию о текущем авторизованном пользователе.

**Успешный ответ (200 OK)** – объект `User`.

### GET `/user/{id}`
Получить информацию о любом пользователе по ID.

### PATCH `/user/password`
Сменить пароль текущего пользователя.

**Тело запроса**:
```json
{
  "old_password": "string", // мин. 6 символов
  "new_password": "string"  // мин. 6 символов
}
```

**Успешный ответ (200 OK)** – `{ "success": true, "message": "password changed", "data": null }`

---

## Преподаватели

### POST `/teachers` (требуется `admin`)
Создать нового преподавателя (пользователь + профиль teacher).

**Тело запроса**: такое же как `CreateUserRequest` (username, full_name, email, password).

**Успешный ответ (200 OK)** – объект `Teacher`.

### POST `/teachers/assign` (требуется `admin`)
Назначить существующего пользователя преподавателем.

**Тело запроса**:
```json
{
  "user_id": 123
}
```

### DELETE `/teachers/{id}` (требуется `admin`)
Удалить профиль преподавателя (пользователь остаётся с ролью `user`).

### GET `/teachers`
Получить список всех преподавателей.

### GET `/teachers/{id}`
Получить информацию о преподавателе по ID пользователя.

---

## Студенты

### POST `/students` (требуется `admin`)
Создать нового студента с привязкой к группе.

**Тело запроса**:
```json
{
  "username": "string",
  "full_name": "string",
  "email": "string",
  "password": "string",
  "group_id": 123
}
```

**Успешный ответ (200 OK)** – объект `Student`.

### POST `/students/assign` (требуется `admin`)
Назначить существующего пользователя студентом и привязать к группе.

**Тело запроса**:
```json
{
  "user_id": 123,
  "group_id": 456
}
```

### PATCH `/students/{id}/group` (требуется `admin`)
Изменить группу студента.

**Тело запроса**:
```json
{
  "group_id": 456
}
```

### DELETE `/students/{id}` (требуется `admin`)
Удалить профиль студента.

### GET `/students`
Получить список всех студентов.

### GET `/students/{id}`
Получить информацию о студенте по ID пользователя.

---

## Группы

### POST `/groups` (требуется `admin`)
Создать учебную группу.

**Тело запроса**:
```json
{
  "name": "ИС-41",
  "specialty": "Информатика",
  "admission_year": 2024
}
```

### GET `/groups/template` (требуется `admin`)
Скачать Excel-шаблон для заполнения данных группы, дисциплин и студентов.

**Ответ**: файл `group_template.xlsx` с тремя листами:
- `Информация о группе` (name, specialty, admission_year)
- `Дисциплины` (title, semester, hours_load, start_date, end_date)
- `Студенты` (full_name, email)

### POST `/groups/upload` (требуется `admin`)
Загрузить заполненный Excel-файл для создания группы, учебного плана и студентов.

**Тело запроса**: `multipart/form-data` с полем `file`.

**Успешный ответ (200 OK)** – структура `GroupWithCurriculumResponse` (аналогично POST `/groups/with-curriculum`).

**Ошибки**:
- `400` – неверный формат файла или отсутствие данных
- `500` – внутренняя ошибка при импорте

**Успешный ответ (200 OK)** – объект `Group`.

### PATCH `/groups/{id}` (требуется `admin`)
Обновить информацию о группе.

**Тело запроса**: поля `name`, `specialty`, `admission_year` (все опциональны).

### DELETE `/groups/{id}` (требуется `admin`)
Удалить группу.

### GET `/groups`
Получить список всех групп.

### GET `/groups/{id}`
Получить группу по ID.

---

## Аудитории

### POST `/audiences` (требуется `admin`)
Создать аудиторию.

**Тело запроса**:
```json
{
  "name": "Аудитория 101",
  "number": "101"
}
```

### PATCH `/audiences/{id}` (требуется `admin`)
Обновить аудиторию.

### DELETE `/audiences/{id}` (требуется `admin`)
Удалить аудиторию.

### GET `/audiences`
Получить список всех аудиторий.

### GET `/audiences/{id}`
Получить аудиторию по ID.

---

## Предметы

### POST `/subjects` (требуется `admin`)
Создать предмет.

**Тело запроса**:
```json
{
  "title": "Математика",
  "semester": 1,
  "hours_load": 72,
  "start_date": "2024-09-01",
  "end_date": "2024-12-31",
  "group_id": 123
}
```

### PATCH `/subjects/{id}` (требуется `admin`)
Обновить предмет.

**Тело запроса**: любые поля из `CreateSubjectRequest` (кроме `group_id`).

### DELETE `/subjects/{id}` (требуется `admin`)
Удалить предмет.

### GET `/subjects`
Получить список предметов с пагинацией.

**Параметры query**: `page`, `per_page`, `sort_by`, `sort_order` (поля: `id`, `title`, `semester`, `hours_load`, `start_date`, `end_date`, `group_id`).

**Успешный ответ**:
```json
{
  "success": true,
  "message": "subjects retrieved",
  "data": {
    "subjects": [ ... ],
    "total": 100,
    "page": 1,
    "per_page": 10,
    "total_pages": 10
  }
}
```

### GET `/subjects/group/{group_id}`
Получить предметы конкретной группы (с пагинацией).

### GET `/subjects/{id}`
Получить предмет по ID.

---

## Учебные периоды (Academic periods)

### POST `/academic-periods` (требуется `admin`)
Создать учебный период (семестр).

**Тело запроса**:
```json
{
  "year": "2024-2025",
  "semester": 1,
  "start_date": "2024-09-01",
  "end_date": "2024-12-31"
}
```

### PATCH `/academic-periods/{id}` (требуется `admin`)
Обновить учебный период.

### DELETE `/academic-periods/{id}` (требуется `admin`)
Удалить учебный период.

### GET `/academic-periods/{id}`
Получить период по ID.

### GET `/academic-periods`
Получить список всех периодов.

---

## Расписание (шаблоны)

Шаблоны расписания привязаны к учебному периоду, группе (через предмет), преподавателю, аудитории.

### POST `/schedule` (требуется `admin`)
Создать шаблон занятия.

**Тело запроса**:
```json
{
  "day_of_week": 1,        // 1-7 (пн-вс)
  "number": 1,             // номер пары (1..8)
  "week_type": 0,          // 0 – обе недели, 1 – числитель, 2 – знаменатель
  "subject_id": 123,
  "teacher_id": 456,
  "audience_id": 789,
  "academic_period_id": 10
}
```

**Успешный ответ (200 OK)** – объект `ScheduleTemplate`.

### PATCH `/schedule/{id}` (требуется `admin`)
Обновить шаблон.

**Тело запроса**: любые поля из `CreateScheduleTemplateRequest` (все опциональны).

### DELETE `/schedule/{id}` (требуется `admin`)
Удалить шаблон.

### POST `/schedule/semester` (требуется `admin`)
Массовое создание расписания на семестр для группы. Проверяет отсутствие конфликтов и автоматически генерирует planned занятия в журнале.

**Тело запроса**:
```json
{
  "group_id": 123,
  "academic_period_id": 10,
  "schedule_entries": [
    {
      "day_of_week": 1,
      "number": 1,
      "week_type": 0,
      "subject_id": 100,
      "teacher_id": 200,
      "audience_id": 300
    },
    ...
  ]
}
```

### GET `/schedule`
Получить список шаблонов с фильтрацией.

**Параметры query**:
- `group_id` (int)
- `subject_id` (int)
- `teacher_id` (int)
- `audience_id` (int)
- `day_of_week` (int)
- `week_type` (int)
- `academic_period_id` (int)

**Успешный ответ (200 OK)** – массив `ScheduleTemplate`.

### GET `/schedule/{id}`
Получить шаблон по ID.

### GET `/schedule/group/{group_id}`
Получить расписание группы (активный период или указанный `?period_id=...`).

### GET `/schedule/teacher/{teacher_id}`
Получить расписание преподавателя.

### GET `/schedule/audience/{audience_id}`
Получить расписание аудитории.

---

## Журнал занятий (Lesson logs)

Журнал содержит фактические занятия (planned, conducted, cancelled, rescheduled). Planned занятия автоматически генерируются из шаблонов расписания с учётом нагрузки предмета.

### GET `/lessons`
Получить записи журнала с фильтрацией.

**Параметры query**:
- `group_id` (int)
- `subject_id` (int)
- `teacher_id` (int)
- `audience_id` (int)
- `academic_period_id` (int)
- `date_from` (string, YYYY-MM-DD)
- `date_to` (string, YYYY-MM-DD)
- `status` (string: `planned`, `completed`, `canceled`, `rescheduled`)
- `number` (int, номер пары)

**Успешный ответ (200 OK)** – массив `LessonLog`.

### POST `/lessons/cancel/{id}` (требуется `admin`)
Отменить занятие по ID (меняет статус на `canceled` и пересчитывает planned занятия предмета).

**Тело запроса**:
```json
{
  "comment": "Причина отмены"
}
```

### POST `/lessons/reschedule` (требуется `admin`)
Создать перенесённое занятие (статус `rescheduled`). Проверяет конфликты с другими занятиями.

**Тело запроса**:
```json
{
  "date": "2024-10-10",
  "number": 2,
  "subject_id": 123,
  "teacher_id": 456,
  "audience_id": 789,
  "comment": "Перенос по болезни"
}
```

### POST `/lessons/complete` (требуется `admin`)
Отметить все занятия за указанную дату как проведённые (`completed`). Занятия со статусами `canceled` или `completed` не изменяются.

**Тело запроса**:
```json
{
  "date": "2024-10-10",
  "comment": "Все занятия проведены"
}
```

---

## Модели данных

### User
```json
{
  "id": 1,
  "username": "john_doe",
  "full_name": "John Doe",
  "email": "john@example.com",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "roles": [
    { "id": 1, "name": "user", "description": "Обычный пользователь" }
  ]
}
```

### Group
```json
{
  "id": 1,
  "name": "ИС-41",
  "specialty": "Информатика",
  "admission_year": 2024
}
```

### Student
```json
{
  "user": { ... },
  "group": { ... }  // может быть null
}
```

### Teacher
```json
{
  "user": { ... }
}
```

### Audience
```json
{
  "id": 1,
  "name": "Аудитория 101",
  "number": "101"
}
```

### Subject
```json
{
  "id": 1,
  "title": "Математика",
  "semester": 1,
  "hours_load": 72,
  "start_date": "2024-09-01T00:00:00Z",
  "end_date": "2024-12-31T00:00:00Z",
  "group_id": 1,
  "group": { ... }
}
```

### AcademicPeriod
```json
{
  "id": 1,
  "year": "2024-2025",
  "semester": 1,
  "start_date": "2024-09-01T00:00:00Z",
  "end_date": "2024-12-31T00:00:00Z",
  "created_at": "...",
  "updated_at": "..."
}
```

### ScheduleTemplate
```json
{
  "id": 1,
  "day_of_week": 1,
  "number": 1,
  "week_type": 0,
  "subject_id": 123,
  "teacher_id": 456,
  "audience_id": 789,
  "academic_period_id": 10,
  "academic_period": { ... },
  "subject": { ... },
  "teacher": { ... },
  "audience": { ... }
}
```

### LessonLog
```json
{
  "id": 1,
  "date": "2024-09-02T00:00:00Z",
  "number": 1,
  "status": "planned",
  "comment": "Сгенерировано автоматически",
  "subject_id": 123,
  "teacher_id": 456,
  "audience_id": 789,
  "academic_period_id": 10,
  "is_from_template": true,
  "academic_period": { ... },
  "subject": { ... },
  "teacher": { ... },
  "audience": { ... }
}
```

### Session
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "refresh_token": "string",
  "user_agent": "Mozilla/5.0 ...",
  "client_ip": "192.168.1.1",
  "created_at": "2024-01-01T00:00:00Z",
  "user": { ... }
}
```

### Role
```json
{
  "id": 1,
  "name": "admin",
  "description": "Администратор"
}
```

---

## Примечания по реализации

- При создании/изменении шаблонов расписания автоматически обновляются planned занятия в журнале (`lesson_logs`) с учётом нагрузки предмета (`hours_load`). Занятия распределяются по датам в пределах `start_date`–`end_date` предмета в соответствии с днями недели и типом недели.
- Отмена занятия (`cancel`) также пересчитывает planned занятия, чтобы сохранить необходимое количество пар.
- При массовом создании расписания на семестр проверяется отсутствие конфликтов (преподаватель, аудитория, группа в одно время).
- Пагинация для списков пользователей и предметов реализована стандартными параметрами `page`, `per_page`, `sort_by`, `sort_order`.

## Дополнения к API документации

### 1. Эндпоинт создания группы с учебным планом и студентами

**POST `/groups/with-curriculum`** (требуется роль `admin`)  
Создаёт учебную группу, список предметов (учебный план) и список студентов с автоматической генерацией логинов и паролей.

**Тело запроса**:
```json
{
  "group": {
    "name": "ИС-51",
    "specialty": "Информационные системы",
    "admission_year": 2025
  },
  "subjects": [
    {
      "title": "Математика",
      "semester": 1,
      "hours_load": 72,
      "start_date": "2025-09-01",
      "end_date": "2025-12-31"
    }
  ],
  "students": [
    {
      "full_name": "Иванов Иван Иванович",
      "email": "ivanov@example.com"
    }
  ]
}
```

**Поля**:
- `group` (object, обязательное) – данные группы (поля `name`, `specialty`, `admission_year`).
- `subjects` (array, обязательное, минимум 1 элемент) – массив предметов, каждый с полями:
  - `title` – название предмета
  - `semester` – номер семестра (1-10)
  - `hours_load` – общее количество часов
  - `start_date` – дата начала изучения (YYYY-MM-DD)
  - `end_date` – дата окончания (YYYY-MM-DD)
- `students` (array, обязательное, минимум 1 элемент) – массив студентов, каждый с полями:
  - `full_name` – ФИО
  - `email` – email (должен быть валидным)

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "group with curriculum and students created",
  "data": {
    "group": { ... },
    "subjects": [ ... ],
    "students": [
      {
        "user": { ... },
        "password": "сгенерированный_пароль",
        "group_id": 123
      }
    ]
  }
}
```

**Примечания**:
- Логин студента генерируется автоматически на основе ФИО и номера группы (например, `501ivanov`).
- Пароль генерируется случайным образом (10 символов) и возвращается в ответе.
- Предметы привязываются к созданной группе.
- Студентам автоматически назначается роль `student` и привязка к группе.

---

### 2. Эндпоинт чат-бота (LLM ассистент)

**POST `/chat`** (требуется аутентификация)  
Отправляет сообщение ассистенту, который умеет отвечать на вопросы о расписании, используя инструменты (например, получение расписания группы). Если запрос не связан с расписанием, отвечает как обычный ассистент.

**Тело запроса**:
```json
{
  "message": "string"   // обязательное, текст запроса пользователя
}
```

**Успешный ответ (200 OK)**:
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "reply": "текст ответа ассистента"
  }
}
```

**Пример запроса**:
```json
{
  "message": "Покажи расписание группы 501 на эту неделю"
}
```

**Примечания**:
- Ассистент сам определяет, нужно ли вызывать инструменты (например, `get_schedule_group`).
- Поддерживаются вопросы о расписании на сегодня, завтра, неделю, с указанием группы и учебного периода.
- При отсутствии необходимых параметров (например, номера группы) ассистент может запросить уточнение.
- Используется LLM-провайдер, заданный в конфигурации (Ollama или GigaChat).

--- 

*Документация актуальна на основе исходного кода сервиса `core`.*
