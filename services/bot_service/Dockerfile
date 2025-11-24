FROM python:3.11-slim

WORKDIR /app

# Установка uv для быстрой установки зависимостей
RUN pip install uv

# Копируем файлы зависимостей
COPY pyproject.toml ./

# Устанавливаем зависимости
RUN uv pip install --system -r pyproject.toml

# Копируем исходный код
COPY . .

# Создаем не-root пользователя
RUN useradd -m -u 1000 app && chown -R app:app /app
USER app

# Экспонируем порт
EXPOSE 8000

# Запуск приложения
CMD ["uvicorn", "app.main:app", "--reload", "--host", "0.0.0.0", "--port", "8000"]
