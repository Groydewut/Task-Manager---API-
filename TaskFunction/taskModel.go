package taskfunction

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
}

type TaskModel struct {
	DB *sql.DB
}

func InitDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

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

	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id          SERIAL PRIMARY KEY,           -- автоинкремент, уникальный ID
		title       VARCHAR(255) NOT NULL,        -- название задачи, обязательное
		description TEXT,                         -- описание, может быть пустым
		status      VARCHAR(20) NOT NULL DEFAULT 'pending',
		priority    VARCHAR(20) NOT NULL DEFAULT 'medium',
		created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать таблицу - %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL")

	return DB, nil
}
