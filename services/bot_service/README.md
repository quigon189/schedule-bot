## Схема взаимодействия пользователя с ботом

```mermide
sequenceDiagram
    participant U as User
    participant TG as Telegram API
    participant B as Bot Instance
    participant DP as Dispatcher
    participant MW as Middleware
    participant R as Router
    participant H as Handler
    participant DB as Database/Service

    Note over U,DB: Пользователь отправляет сообщение
    U->>TG: Отправляет сообщение/команду
    TG->>B: Передаёт Update
    B->>DP: Передаёт Update
    
    Note over DP,MW: Middleware Layer
    DP->>MW: Проверка пользователя
    MW->>DB: Проверка регистрации
    DB-->>MW: Статус пользователя
    alt Не зарегистрирован
        MW-->>B: Блокировка/Приветствие
        B-->>TG: Ответ "Зарегистрируйтесь"
        TG-->>U: Сообщение
    else Зарегистрирован
        MW-->>DP: Пропускает дальше
    end
    
    Note over DP,R: Router Dispatch
    DP->>R: Выбор подходящего Router
    alt Команда (/start, /help)
        R->>CR: Command Router
        CR->>CH: Command Handler
        CH-->>B: Формирует ответ
    else Текстовое сообщение
        R->>ER: Echo Router
        ER->>EH: Echo Handler
        EH-->>B: Формирует ответ
    else Callback от кнопки
        R->>CBR: Callback Router
        CBR->>CBH: Callback Handler
        CBH-->>B: Формирует ответ
    else Другие типы
        R->>Other: Соответствующий Router
        Other->>OH: Другие Handlers
        OH-->>B: Формирует ответ
    end
    
    B->>TG: Отправляет ответ
    TG->>U: Показывает сообщение
```
