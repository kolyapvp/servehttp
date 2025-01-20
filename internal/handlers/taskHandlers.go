package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"pet1/internal/taskService"
	"strconv"
)

type Handler struct {
	Service *taskService.TaskService
}

// Нужна для создания структуры Handler на этапе инициализации приложения

func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Service.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task taskService.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := h.Service.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdTask)
}

// UpdateTaskHandler обрабатывает PATCH-запросы для обновления задачи по ID
func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Преобразуем ID из строки в uint
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Декодируем JSON из тела запроса в структуру Task
	var task taskService.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Вызываем метод сервиса для обновления задачи
	updatedTask, err := h.Service.UpdateTaskByID(uint(id), task)
	if err != nil {
		if err.Error() == "task not found" {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Отправляем обновлённую задачу в ответе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

// DeleteTaskHandler обрабатывает DELETE-запросы для удаления задачи по ID
func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Преобразуем ID из строки в uint
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Вызываем метод сервиса для удаления задачи
	err = h.Service.DeleteTaskByID(uint(id))
	if err != nil {
		if err.Error() == "task not found" || err.Error() == "no task was deleted" {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
