package TaskController

import (
	taskfunction "TaskManager-API/TaskFunction"
	"context"
	"net/http"
	"time"
)

type Handler struct {
	TaskM taskfunction.TaskModel
}

func (h Handler) CreateTask(w http.ResponseWriter, r *http.Request) {}

func (h Handler) TaskList(w http.ResponseWriter, r *http.Request) {}

func (h Handler) TaskListByID(w http.ResponseWriter, r *http.Request) {}

func (h Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {}

func (h Handler) DeleateTask(w http.ResponseWriter, r *http.Request) {}

func (h Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {

	db, _ := taskfunction.InitDB()
	defer db.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("data base is down"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
