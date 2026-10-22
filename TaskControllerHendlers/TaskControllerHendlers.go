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
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	TaskM taskfunction.TaskModel
	Redis *redis.Client
}

const cacheTTL = 5 * time.Minute

func (h Handler) setCache(ctx context.Context, key string, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	h.Redis.Set(ctx, key, data, cacheTTL)
}

func (h Handler) getCache(ctx context.Context, key string, dest any) bool {
	val, err := h.Redis.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false
	}
	return true
}

func (h Handler) invalidateCache(ctx context.Context, taskID int) {
	h.Redis.Del(ctx, fmt.Sprintf("tasks:%d", taskID))
	iter := h.Redis.Scan(ctx, 0, "tasks:list:*", 0).Iterator()
	for iter.Next(ctx) {
		h.Redis.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		slog.Error("ошибка сканирования ключей Redis", "err", err)
	}
}

func (h Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req taskfunction.Task
	ctx := context.Background()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if err := taskfunction.ValidateTask(req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	requestID := r.Header.Get("X_Request_ID")
	if requestID == "" {
		slog.Error(`{"error":"X-Request-ID header is required"}`)
		return
	}
	lockKey := fmt.Sprintf("request:%s", requestID)
	ok, err := h.Redis.SetNX(ctx, lockKey, "processing", time.Second*60).Result()
	if err != nil {
		slog.Warn("redis недоступен", "err", err)
	} else if !ok {
		http.Error(w, `{"error":"duplicate request"}`, http.StatusConflict)
		return
	}
	id, err := h.TaskM.InsertTask(req)
	if err != nil {
		slog.Error("create task failed", "error", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	// Инвалидация: списки изменились
	h.invalidateCache(r.Context(), 0) // 0 — не удаляем конкретную задачу, только списки

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":      id,
		"message": "task created",
	})

}

func (h Handler) TaskList(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	cacheKey := fmt.Sprintf("tasks:list:%s:%d:%d", status, limit, offset)
	var tasks []taskfunction.Task
	if h.getCache(r.Context(), cacheKey, &tasks) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(tasks)
		return
	}

	tasks, err := h.TaskM.GetAllTasks(status, limit, offset)
	if err != nil {
		slog.Error("get tasks error", "err", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	h.setCache(r.Context(), cacheKey, tasks)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(tasks)
}

func (h Handler) TaskListByID(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if id <= 0 {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	cacheKey := fmt.Sprintf("task:%d", id)

	var task taskfunction.Task
	if h.getCache(r.Context(), cacheKey, &task) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(task)
		return
	}
	task, err := h.TaskM.GetByID(id)
	if err != nil {
		if errors.Is(err, taskfunction.ErrTaskNotFound) {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	h.setCache(r.Context(), cacheKey, task)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(task)
}

func (h Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		http.Error(w, "Отправлены не верные данные", http.StatusBadRequest)
		return
	}

	var req taskfunction.Task

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := taskfunction.ValidateTask(req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := h.TaskM.Update(req); err != nil {
		if errors.Is(err, taskfunction.ErrTaskNotFound) {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	h.invalidateCache(r.Context(), id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "task updated"})
}

func (h Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
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

	h.invalidateCache(r.Context(), id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Элемент успешно удалён."}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode delete response", "err", err)
	}

}

func (h Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.TaskM.DB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "db": err.Error()})
		return
	}

	// Проверим и Redis
	if err := h.Redis.Ping(ctx).Err(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "redis": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})

}
