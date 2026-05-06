.PHONY: help dev up down logs build clean

# Файлы docker-compose для обычного и dev-режима
COMPOSE_FILES = -f compose.yaml -f docker-compose.dev.yaml

help:
	@echo "Доступные команды:"
	@echo "  make dev   - запустить все сервисы в режиме разработки (core с air)"
	@echo "  make up    - запустить все сервисы в обычном режиме"
	@echo "  make down  - остановить все сервисы"
	@echo "  make logs  - показать логи всех сервисов (Ctrl+C для выхода)"
	@echo "  make build - пересобрать образы"
	@echo "  make clean - остановить и удалить контейнеры + volumes"

dev:
	docker-compose -f docker-compose.dev.yaml up -d
	@echo "✓ Сервисы запущены в режиме разработки. Логи core: docker-compose logs -f core"

up:
	docker-compose -f compose.yaml up -d

down:
	docker-compose $(COMPOSE_FILES) down

logs:
	docker-compose $(COMPOSE_FILES) logs -f

build:
	docker-compose $(COMPOSE_FILES) build --no-cache

clean:
	docker-compose $(COMPOSE_FILES) down -v
