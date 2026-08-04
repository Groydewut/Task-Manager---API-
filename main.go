package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	middleware "TaskManager-API/Middleware"
	TaskController "TaskManager-API/TaskControllerHendlers"
	taskfunction "TaskManager-API/TaskFunction"
)

func main() {
	//! Регистрация логера

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Не удалось загрузить файл .env, используются переменные окружения")
		os.Exit(1)
	}

	//! Подключение к бд
	db, err := taskfunction.InitDB()
	if err != nil {
		logger.Error("Ошибка подключения к БД", slog.String("err", err.Error()))
		os.Exit(1)
	}
	logger.Info("Успешное подключение к PostgreSQL")

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(8)

	defer func() {
		logger.Info("Закрытие пула соеденений базы данных")
		db.Close()
	}()

	if db == nil {
		logger.Error("Критическая ошибка. Переменая базы данных равна nil")
		os.Exit(1)
	}

	myModel := taskfunction.TaskModel{DB: db}
	myHandler := TaskController.Handler{
		TaskM: myModel,
	}

	//! Создание роутера
	r := chi.NewRouter()

	//! Маршруты
	r.Group(func(r chi.Router) {
		r.Use(middleware.SlogMiddleware(logger))

		r.Post("/tasks", myHandler.CreateTask)
		r.Get("/tasks", myHandler.TaskList)
		r.Get("/tasks/{id}", myHandler.TaskListByID)
		r.Put("/tasks/{id}", myHandler.UpdateTask)
		r.Delete("/tasks/{id}", myHandler.DeleteTask)
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
		logger.Info("Сервер запущен", slog.String("url", "http://localhost:8080"))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Ошибка запуска сервера", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	//! Блокируем поток до получения сигнала
	sig := <-shutdownSignal
	logger.Info("Получен системный сигнал. Начинаем Graceful Shutdown...", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//! Останавливаем сервер, он перестаёт проинимать новые запросы и ждет обработки старых в пределах таймаута
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Ошибка при остановке сервера", slog.String("err", err.Error()))
		os.Exit(1)
	}
	logger.Info("Сервер успешно остановлен")
}
