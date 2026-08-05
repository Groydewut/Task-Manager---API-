# Task Manager API

REST API для управления задачами на **Go** с веб-интерфейсом.
Проект демонстрирует production-ready подход: контейнеризация, миграции БД, graceful shutdown, структурированное логирование.

---

## 🚀 Стек технологий

| Категория       | Технологии                                                                 |
|-----------------|----------------------------------------------------------------------------|
| **Язык**        | Go 1.26.2                                                                  |
| **Роутер**      | [chi/v5](https://github.com/go-chi/chi)                                    |
| **База данных** | PostgreSQL 17 (`database/sql` + драйвер `lib/pq`)                          |
| **Миграции**    | [golang-migrate](https://github.com/golang-migrate/migrate)                |
| **Конфиг**      | Переменные окружения + [godotenv](https://github.com/joho/godotenv)        |
| **Логирование** | [`log/slog`](https://pkg.go.dev/log/slog) (структурированные логи)         |
| **Контейнеры**  | Docker (multi-stage build) + Docker Compose                                |
| **Frontend**    | HTML5 + CSS3 + Vanilla JavaScript                                          |

---

## 📋 Возможности

### API Endpoints

| Метод  | Путь           | Описание                               |
| ------ | -------------- | -------------------------------------- |
| POST   | `/tasks`       | Создать новую задачу                   |
| GET    | `/tasks`       | Список задач (с фильтрацией и пагинацией) |
| GET    | `/tasks/{id}`  | Получить задачу по ID                  |
| PUT    | `/tasks/{id}`  | Обновить задачу                        |
| DELETE | `/tasks/{id}`  | Удалить задачу                         |
| GET    | `/health`      | Проверка здоровья (БД + сервер)        |

### Функционал

- ✅ **CRUD операции** над задачами
- ✅ **Фильтрация** по статусу (`status=pending|in_progress|done`)
- ✅ **Пагинация** (`limit`, `offset`)
- ✅ **Валидация** входных данных (title, status, priority, description)
- ✅ **Graceful shutdown** — корректная остановка сервера
- ✅ **Auto-migrations** — применение миграций при старте
- ✅ **Connection pooling** — настройка пула соединений с БД
- ✅ **Структурированные логи** с метриками запросов
- ✅ **Веб-интерфейс** — создание, редактирование, удаление задач

---

## 🗂 Модель данных

### Task

```json
{
  "id": 1,
  "title": "string (required, max 255)",
  "description": "string (optional, max 5000)",
  "status": "pending|in_progress|done",
  "priority": "low|medium|high",
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:00:00Z"
}
```

### Схема БД

```sql
CREATE TABLE tasks (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      VARCHAR(20) DEFAULT 'pending',
    priority    VARCHAR(20) DEFAULT 'medium',
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## ⚙️ Конфигурация

Переменные окружения:

| Переменная      | Описание                    | По умолчанию |
|-----------------|-----------------------------|--------------|
| `DB_HOST`       | Хост PostgreSQL             | `localhost`  |
| `DB_PORT`       | Порт PostgreSQL             | `5432`       |
| `DB_USER`       | Пользователь БД             | `postgres`   |
| `DB_PASSWORD`   | Пароль БД                   | `secret`     |
| `DB_NAME`       | Имя базы данных             | `TaskDb`     |
| `PORT`          | Порт HTTP-сервера           | `8080`       |

---

## 📦 Запуск

### Вариант 1: Docker Compose (рекомендуется)

Одной командой поднимает БД и API:

```bash
docker compose up --build
```

Сервис будет доступен на: **http://localhost:8080**

### Вариант 2: Локальная разработка

#### 1. Клонировать репозиторий

```bash
git clone https://github.com/Groydewut/Task-Manager---API-.git
cd Task-Manager---API-
```

#### 2. Создать `.env` файл

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=TaskDb
```

#### 3. Запустить PostgreSQL

```bash
# Или используйте docker compose up db
docker run -d \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=TaskDb \
  -p 5432:5432 \
  postgres:17
```

#### 4. Применить миграции

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path migrations -database "postgres://postgres:secret@localhost:5432/TaskDb?sslmode=disable" up
```

#### 5. Запустить сервер

```bash
go mod download
go run main.go
```

---

## 🏗 Структура проекта

```
.
├── main.go                     # Точка входа, роутинг, graceful shutdown
├── go.mod                      # Зависимости Go
├── docker-compose.yml          # Оркестрация контейнеров
├── Dockerfile.v3               # Multi-stage Dockerfile
├── Middleware/
│   └── middleware.go           # Slog-логгер для запросов
├── TaskControllerHendlers/
│   └── TaskControllerHendlers.go # HTTP-обработчики (handlers)
├── TaskFunction/
│   └── taskModel.go            # Модель Task, работа с БД, валидация
├── frontend/
│   ├── index.html              # UI разметка
│   ├── styles.css              # Стили
│   ├── app.js                  # Клиентская логика
│   └── favicon.png             # Иконка
└── migrations/
    ├── 000001_create_tasks_table.up.sql
    └── 000001_create_tasks_table.down.sql
```

---

## 🔧 Особенности реализации

### Graceful Shutdown
Сервер корректно обрабатывает сигналы `SIGINT`/`SIGTERM`:
- Останавливает приём новых запросов
- Даёт завершиться текущим запросам (до 5 секунд)
- Закрывает соединения с БД

### Connection Pooling
```go
db.SetMaxOpenConns(15)  // Максимум открытых соединений
db.SetMaxIdleConns(8)   // Минимум в пуле
```

### Auto-Migrations
Миграции применяются автоматически при подключении к БД через `golang-migrate`.

### Логирование
Каждый запрос логируется с метриками:
```
time=... level=INFO msg="incoming request" method=GET path=/tasks status=200 duration_ms=5
```

### Безопасность
- Non-root пользователь в Docker-контейнере
- Валидация всех входных данных
- SQL-инъекции предотвращены (prepared statements)

---

## 🧪 Примеры запросов

### Создать задачу

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Изучить Go",
    "description": "Прочитать документацию",
    "status": "pending",
    "priority": "high"
  }'
```

### Получить список задач

```bash
curl http://localhost:8080/tasks?status=pending&limit=10&offset=0
```

### Обновить задачу

```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Изучить Go",
    "status": "in_progress",
    "priority": "high"
  }'
```

### Healthcheck

```bash
curl http://localhost:8080/health
# Ответ: OK (статус 200)
```

---

## 📊 Статус проекта

| Этап              | Статус     |
|-------------------|------------|
| CRUD API          | ✅ Готово  |
| Frontend UI       | ✅ Готово  |
| Graceful Shutdown | ✅ Готово  |
| Миграции БД       | ✅ Готово  |
| Docker            | ✅ Готово  |
| Логирование       | ✅ Готово  |
| Валидация         | ✅ Готово  |

---

## 🔮 Планы развития

- [ ] Unit-тесты (mocks для handlers/storage)
- [ ] Интеграционные тесты
- [ ] CI/CD пайплайн (GitHub Actions)
- [ ] Rate limiting
- [ ] Аутентификация/авторизация (JWT)
- [ ] Метрики Prometheus (`/metrics`)
- [ ] OpenAPI/Swagger документация

---

## 📝 Лицензия

MIT

---

## 👤 Автор

Учебный проект для отработки практик докеризации, CI/CD и production-разработки.
