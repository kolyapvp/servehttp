package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pet1/internal/taskService"
	"pet1/internal/web/tasks"
	"strconv"
)

type Handler struct {
	Service *taskService.TaskService
}

// DeleteTasksId реализует удаление задачи по ID
func (h *Handler) DeleteTasksId(ctx context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	// Извлекаем ID задачи из запроса
	id := request.Id

	// Вызываем сервис для удаления задачи
	err := h.Service.DeleteTaskByID(id)
	if err != nil {
		if err.Error() == "task not found" || err.Error() == "no task was deleted" {
			// Возвращаем 404 Not Found, если задача не найдена
			return tasks.DeleteTasksId404Response{}, nil
		}
		// Возвращаем 500 Internal Server Error для других ошибок
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}

	// Возвращаем 204 No Content при успешном удалении
	return tasks.DeleteTasksId204Response{}, nil
}

func (h *Handler) PatchTasksId(_ context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	// Извлекаем ID задачи из запроса
	id := request.Id

	// Извлекаем тело запроса, содержащее поля для обновления
	body := request.Body

	// Создаём объект задачи с обновлёнными полями
	taskToUpdate := taskService.Task{}

	if body.Task != nil {
		taskToUpdate.Task = *body.Task
	}

	if body.IsDone != nil {
		taskToUpdate.IsDone = *body.IsDone
	}

	// Вызываем сервис для обновления задачи
	updatedTask, err := h.Service.UpdateTaskByID(id, taskToUpdate)
	if err != nil {
		if err.Error() == "task not found" {
			// Возвращаем 404 Not Found, если задача не найдена
			return tasks.PatchTasksId404Response{}, nil
		}
		// Возвращаем 500 Internal Server Error для других ошибок
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Создаём объект ответа с обновлённой задачей
	responseTask := tasks.Task{
		Id:     &updatedTask.ID,
		Task:   &updatedTask.Task,
		IsDone: &updatedTask.IsDone,
	}

	// Возвращаем 200 OK с обновлённой задачей
	return tasks.PatchTasksId200JSONResponse(responseTask), nil
}

func (h *Handler) GetTasks(_ context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	// Получение всех задач из сервиса
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	// Создаем переменную респон типа 200джейсонРеспонс
	// Которую мы потом передадим в качестве ответа
	response := tasks.GetTasks200JSONResponse{}

	// Заполняем слайс response всеми задачами из БД
	for _, tsk := range allTasks {
		task := tasks.Task{
			Id:     &tsk.ID,
			Task:   &tsk.Task,
			IsDone: &tsk.IsDone,
		}
		response = append(response, task)
	}

	// САМОЕ ПРЕКРАСНОЕ. Возвращаем просто респонс и nil!
	return response, nil
}

func (h *Handler) PostTasks(_ context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	// Распаковываем тело запроса напрямую, без декодера!
	taskRequest := request.Body
	// Обращаемся к сервису и создаем задачу
	taskToCreate := taskService.Task{
		Task:   *taskRequest.Task,
		IsDone: *taskRequest.IsDone,
	}
	createdTask, err := h.Service.CreateTask(taskToCreate)

	if err != nil {
		return nil, err
	}
	// создаем структуру респонс
	response := tasks.PostTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   &createdTask.Task,
		IsDone: &createdTask.IsDone,
	}
	// Просто возвращаем респонс!
	return response, nil
}

// Нужна для создания структуры Handler на этапе инициализации приложения

func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
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
