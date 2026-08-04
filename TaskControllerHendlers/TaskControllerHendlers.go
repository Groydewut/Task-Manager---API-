package TaskController

import (
	taskfunction "TaskManager-API/TaskFunction"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	TaskM taskfunction.TaskModel
}

func (h Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req taskfunction.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if err := taskfunction.ValidateTask(req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	id, err := h.TaskM.InsertTask(req)
	if err != nil {
		slog.Error("create task failed", "error", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":      id,
		"message": "task created",
	})

}

func (h Handler) TaskList(w http.ResponseWriter, r *http.Request) {}

func (h Handler) TaskListByID(w http.ResponseWriter, r *http.Request) {}

func (h Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {}

func (h Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Отправлены не верные данные", http.StatusBadRequest)
		return
	}

	if err := h.TaskM.Delete(id); err != nil {
		if errors.Is(err, taskfunction.ErrTaskNotFound) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		slog.Error("delete task failed", "error", err, "task_id", id)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Элемент успешно удалён."}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusBadRequest)
		return
	}

}

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
