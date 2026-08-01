# Task Manager API

REST API для управления задачами, написанный на **Go**.  
Учебный проект «подопытный» — на нём отрабатываются этапы докеризации, CI/CD и продакшн-практики.

---

## 🚀 Стек

- **Go** 1.22+
- **Роутер:** [chi](https://github.com/go-chi/chi) (`v5`)
- **БД:** PostgreSQL (драйвер в процессе подключения)
- **Доступ к БД:** `database/sql` + `sqlx` (или `pgx`) — без ORM
- **Миграции:** [golang-migrate/migrate](https://github.com/golang-migrate/migrate) или [pressly/goose](https://github.com/pressly/goose)
- **Логи:** стандартный [`log/slog`](https://pkg.go.dev/log/slog)
- **Конфиг:** переменные окружения + [godotenv](https://github.com/joho/godotenv)
- **Контейнеризация:** Docker + docker-compose (планируется)

---

## 📋 Функционал (endpoints)

| Метод  | Путь           | Описание                                       |
| ------ | -------------- | ---------------------------------------------- |
| POST   | `/tasks`       | Создать задачу                                 |
| GET    | `/tasks`       | Список задач (фильтр `status`, пагинация)      |
| GET    | `/tasks/{id}`  | Получить задачу по ID                          |
| PUT    | `/tasks/{id}`  | Обновить задачу                                |
| DELETE | `/tasks/{id}`  | Удалить задачу                                 |
| GET    | `/health`      | Healthcheck (проверка соединения с БД)         |

---

## 🗂 Модель `Task`

```go
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title" validate:"required"`
    Description string    `json:"description"`
    Status      string    `json:"status"`     // pending, in_progress, done
    Priority    string    `json:"priority"`   // low, medium, high
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## ⚙️ Конфигурация

Сервис читает переменные окружения:

| Переменная     | Описание                                  | Пример                              |
| -------------- | ----------------------------------------- | ----------------------------------- |
| `PORT`         | Порт HTTP-сервера                         | `8080`                              |
| `DATABASE_URL` | DSN для подключения к PostgreSQL          | `postgres://user:pass@localhost:5432/taskmanager?sslmode=disable` |
| `LOG_LEVEL`    | Уровень логирования (`debug`/`info`/`warn`/`error`) | `info`                     |

Можно положить в `.env` в корне проекта (godotenv подхватит автоматически).

---

## 🛠 Требования к сервису

- ✅ **Graceful shutdown** — корректная обработка `SIGINT`/`SIGTERM`, завершение текущих запросов и закрытие соединений
- ✅ **Конфигурация через env** — `PORT`, `DATABASE_URL`, `LOG_LEVEL`
- ✅ **Миграции** при старте (или отдельной командой)
- ✅ **Структурированные логи** — каждый запрос логируется: метод, путь, статус, время
- ✅ **Валидация** — `title` не пустой, `status`/`priority` только из разрешённых значений

---

## 📦 Запуск

### 🔮 После докеризации (целевой вариант)

Когда проект будет упакован в контейнеры, запуск сервиса и базы данных будет сводиться к **одной команде**:

```bash
docker compose up --build
```

Эта команда поднимет и сам API, и PostgreSQL, и применит миграции — больше ничего не нужно.

### 🧑‍💻 Локальный запуск (пока без Docker)

#### 1. Клонировать репозиторий

```bash
git clone https://github.com/Groydewut/Task-Manager---API-.git
cd Task-Manager---API-
```

#### 2. Подготовить `.env`

```env
PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/taskmanager?sslmode=disable
LOG_LEVEL=info
```

#### 3. Запустить миграции

```bash
# через migrate
migrate -path migrations -database "$DATABASE_URL" up
```

#### 4. Запустить сервис

```bash
go run ./cmd/api
```

Сервис поднимется на `http://localhost:8080`.

---

## 🗺 Планы развития

- [ ] Докеризация (Docker + docker-compose) → запуск одной командой `docker compose up --build`
- [ ] CI/CD (GitHub Actions)
- [ ] Юнит- и интеграционные тесты
- [ ] Метрики (`/metrics`) и трейсинг
- [ ] Авторизация и аутентификация

---

## 📝 Лицензия

MIT (или укажи свою).
