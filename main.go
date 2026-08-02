package main

import (
	"fmt"
	"log"
	"net/http"

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
	defer db.Close()

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

	fmt.Println("Сервер запущен на http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", r))

}
