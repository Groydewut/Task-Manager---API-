package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"TaskManager-API/TaskController"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить файл .env, используются переменные окружения")
	}

	//! Подключение к бд

	//! Создание роутера
	r := chi.NewRouter()

	//! Маршруты
	r.Group(func(r chi.Router) {
		r.Post("/tasks", TaskController.CreateTask)
		r.Get("/tasks", TaskController.TaskList)
		r.Get("/tasks/{id}", TaskController.TaskListByID)
		r.Put("/tasks/{id}", TaskController.UpdateTask)
		r.Delete("/tasks/{id}", TaskController.DeleateTask)
	})
	log.Fatal(http.ListenAndServe(":8080", r))

}
