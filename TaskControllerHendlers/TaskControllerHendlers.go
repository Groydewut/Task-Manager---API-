package TaskController

import (
	taskfunction "TaskManager-API/TaskFunction"
	"net/http"
)

type Handler struct {
	TaskM taskfunction.TaskModel
}

func (h Handler) CreateTask(w http.ResponseWriter, r *http.Request) {}

func (h Handler) TaskList(w http.ResponseWriter, r *http.Request) {}

func (h Handler) TaskListByID(w http.ResponseWriter, r *http.Request) {}

func (h Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {}

func (h Handler) DeleateTask(w http.ResponseWriter, r *http.Request) {}
