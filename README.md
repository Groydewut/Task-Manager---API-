# Task Manager API

REST API для управления задачами, написанный на **Go**.  
Учебный проект «подопытный» — на нём отрабатываются этапы докеризации, CI/CD и продакшн-практики.

---

## 🚀 Стек

- **Go** 1.22+
- **Роутер:** [chi](https://github.com/go-chi/chi) (`v5`)
- **БД:** PostgreSQL (драйвер `database/sql` — стандартная библиотека)
- **Коннект-пул:** настроен для эффективного управления соединениями
- **Миграции:** планируется `golang-migrate` или `goose`
- **Логи:** стандартный [`log/slog`](https://pkg.go.dev/log/slog) + кастомный логгер с единым форматом
- **Конфиг:** переменные окружения + [godotenv](https://github.com/joho/godotenv)
- **Graceful shutdown:** `signal.NotifyContext` + `srv.Shutdown`
- **Контейнеризация:** Docker + docker-compose (планируется)

---

## 📋 Функционал (endpoints)

### ✅ Реализовано

| Метод  | Путь           | Описание                                       | Статус |
| ------ | -------------- | ---------------------------------------------- | ------ |
| POST   | `/tasks`       | Создать задачу                                 | ✅ Готово |
| GET    | `/tasks`       | Список задач (фильтр `status`, пагинация)      | ✅ Готово |
| GET    | `/tasks/{id}`  | Получить задачу по ID                          | ✅ Готово |
| PUT    | `/tasks/{id}`  | Обновить задачу                                | ✅ Готово |
| DELETE | `/tasks/{id}`  | Удалить задачу                                 | ✅ Готово |
| GET    | `/health`      | Healthcheck (проверка соединения с БД)         | ✅ Готово |

### 🛠 Инфраструктура

- ✅ **Graceful shutdown** — корректная обработка `SIGINT`/`SIGTERM`, завершение текущих запросов и закрытие соединений
- ✅ **Конфигурация через env** — `PORT`, `DATABASE_URL`, `LOG_LEVEL`
- ✅ **Коннект-пул** — эффективное управление соединениями с БД
- ✅ **Структурированные логи** — кастомный slog-логгер, единый формат для всех сообщений
- ✅ **Валидация** — `title`, `status`, `description`, `priority` проверяются на корректность

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
# через migrate (когда будет настроено)
migrate -path migrations -database "$DATABASE_URL" up
```

#### 4. Запустить сервис

```bash
go run main.go
```

Сервис поднимется на `http://localhost:8080`.

---

## 🗺 Планы развития

### Этап 1 — MVP (✅ почти завершён)

- [x] Роутер `chi` + регистрация эндпоинтов
- [x] Модель `Task` с JSON-тегами
- [x] Middleware (логирование, обвязка)
- [x] Подключение к PostgreSQL через `database/sql`
- [x] Коннект-пул для БД
- [x] Конфиг из env (`DATABASE_URL`, `PORT`, `LOG_LEVEL`) + `godotenv`
- [x] `slog` — структурное логирование, кастомный логгер
- [x] Graceful shutdown — `signal.NotifyContext` + `srv.Shutdown`
- [x] `GET /health` — healthcheck с проверкой БД
- [x] CRUD по `/tasks` (POST, GET, GET/{id}, PUT, DELETE)
- [x] Валидация полей (`title`, `status`, `description`, `priority`)
- [ ] Миграции — таблица `tasks` через `golang-migrate` или `goose`
- [ ] Repository-слой — вынос SQL в отдельный пакет `internal/storage`

### Этап 2 — Докеризация

- [ ] `Dockerfile` с multi-stage build
- [ ] `docker-compose.yml` с PostgreSQL + миграциями
- [ ] Цель: `docker compose up --build` поднимает всё одной командой

### Этап 3 — Тесты и CI/CD

- [ ] Unit-тесты на хендлеры с моками storage
- [ ] Интеграционные тесты на БД
- [ ] CI/CD через GitHub Actions

### Этап 4 — Продвинутые фичи

- [ ] Метрики `/metrics` и трейсинг
- [ ] Авторизация и аутентификация
- [ ] Rate limiting

---

## 📊 Текущий статус

**Этап 1 — финальная стадия** (готовность ~90% от ТЗ)

✅ **Что работает:**
- Все 5 CRUD-эндпоинтов (`POST`, `GET`, `GET/{id}`, `PUT`, `DELETE`)
- Healthcheck `/health`
- Graceful shutdown
- Структурное логирование через `slog`
- Коннект-пул для БД
- Валидация полей

🟡 **Что осталось:**
- Миграции (таблица `tasks`)
- Repository-слой (вынос SQL в отдельный пакет)

⏳ **Следующий шаг:**
- Закрыть миграции → перейти к Docker

---

## 📝 Лицензия

MIT (или укажи свою).
