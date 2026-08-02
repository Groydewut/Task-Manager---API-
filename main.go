package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	TaskController "TaskManager-API/TaskControllerHendlers"
	taskfunction "TaskManager-API/TaskFunction"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить файл .env, используются переменные окружения")
	}

	//! Подключение к бд
	db, err := taskfunction.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(8)

	defer func() {
		log.Println("Закрытие пула соеденений базы данных")
		db.Close()
	}()

	if db == nil {
		log.Fatal("Критическая ошибка. Переменая базы данных равна nil")
	}

	myModel := taskfunction.TaskModel{DB: db}
	myHandler := TaskController.Handler{
		TaskM: myModel,
	}

	//! Создание роутера
	r := chi.NewRouter()

	//! Маршруты
	r.Group(func(r chi.Router) {
		r.Post("/tasks", myHandler.CreateTask)
		r.Get("/tasks", myHandler.TaskList)
		r.Get("/tasks/{id}", myHandler.TaskListByID)
		r.Put("/tasks/{id}", myHandler.UpdateTask)
		r.Delete("/tasks/{id}", myHandler.DeleateTask)
	})

	r.Get("/health", myHandler.CheckHealth)

	//! Настройка http.Server чтобы иметь доступ к методу Shutdown

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	//! Отслеживание системных сигналов
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	//! Запуск сервера в фоновом режиме
	go func() {
		fmt.Println("Сервер запущен на http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Ошибка запуска сервера")
		}
	}()

	//! Блокируем поток до получения сигнала
	sig := <-shutdownSignal
	log.Printf("Получен сигнал %v. Начинаем Graceful Shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//! Останавливаем сервер, он перестаёт проинимать новые запросы и ждет обработки старых в пределах таймаута
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("ОШибка при остановке сервера: %v", err)
	}
	log.Println("Сервер успешно остановлен")
}
