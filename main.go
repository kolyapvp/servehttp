package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"pet1/db"
	"strconv"
)

// Task представляет структуру задачи
type Task struct {
	ID     uint   `json:"id"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

// POST handler для добавления задачи
func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody Task

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Сохраняем задачу в базу данных
	if err := db.DB.Create(&requestBody).Error; err != nil {
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(requestBody); err != nil {
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
	}
}

// GET handler для получения всех задач
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	// Извлекаем все задачи из базы данных
	if err := db.DB.Find(&tasks).Error; err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
	}
}

// PATCH handler для обновления задачи по ID
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Проверяем, что ID является числом
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Создаём переменную для обновления
	var updatedData Task

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Находим существующую задачу
	var existingTask Task
	if err := db.DB.First(&existingTask, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		return
	}

	// Обновляем поля, если они были предоставлены
	if updatedData.Task != "" {
		existingTask.Task = updatedData.Task
	}
	existingTask.IsDone = updatedData.IsDone

	// Сохраняем обновления в базу данных
	if err := db.DB.Save(&existingTask).Error; err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Отправляем обновлённую задачу в ответе
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingTask)
}

// DELETE handler для удаления задачи по ID
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Проверяем, что ID является числом
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Находим и удаляем задачу
	if err := db.DB.Delete(&Task{}, taskID).Error; err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// Отправляем подтверждение удаления
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Task deleted successfully",
	})
}

func main() {
	// Инициализация базы данных
	db.InitDB()

	// Маршруты
	r := mux.NewRouter()
	r.HandleFunc("/tasks", postTaskHandler).Methods("POST")
	r.HandleFunc("/tasks", getTaskHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", updateTaskHandler).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", deleteTaskHandler).Methods("DELETE")

	// Запускаем сервер
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
