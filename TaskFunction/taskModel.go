package taskfunction

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskModel struct {
	DB *sql.DB
}

func (m TaskModel) GetByID(id int) (Task, error) {
	query := "SELECT id, title, description,status,priority,created_at,updated_at FROM tasks WHERE id = &1"

	var t Task

	err := m.DB.QueryRow(query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status,
		&t.Priority, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return Task{}, fmt.Errorf("задача не найдена")
	}

	if err != nil {
		return Task{}, fmt.Errorf("ошибка при получении задачи: %w", err)
	}
	return t, nil
}

func (m TaskModel) InsertTask(t Task) (int, error) {
	query := "INSERT INTO tasks (title,description,status,priority) VALUES ($1,$2,$3,$4)  RETURNING id"
	var id int

	err := m.DB.QueryRow(query, t.Title, t.Description, t.Status, t.Priority).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("произошла ошибка при добавления записи: %w ", err)
	}
	return id, nil
}

func InitDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	var DB *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Попытка %d из 5 : база данных еще не готова, ждём...", i)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("Не удалось подключиться к бд после 5 попыток: %v", err)
	}

	m, err := migrate.New(
		"file://migrations",
		connStr,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации миграции: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("ошибка применений миграций: %w", err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Info("миграции успешно применены")

	return DB, nil
}

var validateStatuses = map[string]bool{
	"pending":     true,
	"in_progress": true,
	"done":        true,
}

var validPriorities = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
}

func ValidateTask(t Task) error {

	if len(strings.TrimSpace(t.Title)) > 255 {
		return errors.New("title too long")
	}
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(t.Status) != "" && !validateStatuses[t.Status] {
		return errors.New("invalid status")
	}
	if strings.TrimSpace(t.Priority) != "" && !validPriorities[t.Priority] {
		return errors.New("invalid priority")
	}
	if len(strings.TrimSpace(t.Description)) > 5000 {
		return errors.New("description too long")
	}
	return nil
}
