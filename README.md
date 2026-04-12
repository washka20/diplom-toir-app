# ТОиР — Система управления техническим обслуживанием и ремонтом

Веб-приложение для управления техническим обслуживанием и ремонтом оборудования предприятия.

## Стек технологий

- **Backend**: Go 1.22, Echo v4, GORM, PostgreSQL 16
- **Frontend**: Vue 3, TypeScript, Pinia, Vue Router, Vite
- **Инфраструктура**: Docker, Docker Compose, Nginx

## Быстрый старт

```bash
# Клонировать и запустить
git clone <repo>
cd diplom-toir-app
cp .env.example .env
docker compose up --build -d

# Заполнить тестовыми данными
docker compose exec api ./seed
```

Приложение доступно: http://localhost

### Тестовые учётные записи

| Логин | Пароль | Роль |
|-------|--------|------|
| admin | password123 | Администратор |
| engineer | password123 | Инженер |
| technician | password123 | Техник |
| operator | password123 | Оператор |

## Архитектура

```
├── backend/          # Go API (Clean Architecture)
│   ├── cmd/server/   # Точка входа
│   ├── cmd/seed/     # Заполнение тестовыми данными
│   ├── internal/     # Бизнес-логика
│   │   ├── handlers/ # HTTP-обработчики
│   │   ├── services/ # Сервисный слой
│   │   ├── repository/ # Слой данных
│   │   ├── models/   # GORM-модели
│   │   └── middleware/ # JWT, RBAC
│   └── migrations/   # SQL-миграции
├── frontend/         # Vue 3 SPA
│   ├── src/views/    # Страницы
│   ├── src/stores/   # Pinia-хранилища
│   └── src/api/      # HTTP-клиент
└── docker-compose.yml
```

## Тестирование

```bash
cd backend && go test -race -cover ./...
```

## API

Документация API в формате REST. Все эндпоинты возвращают JSON envelope:
```json
{"success": true, "data": {}, "error": null, "meta": {}}
```
