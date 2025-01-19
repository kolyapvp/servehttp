package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pet1/db"
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
	json.NewEncoder(w).Encode(requestBody)
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
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	// Инициализация базы данных
	db.InitDB()

	// Создаем маршруты
	r := mux.NewRouter()
	r.HandleFunc("/tasks", postTaskHandler).Methods("POST")
	r.HandleFunc("/tasks", getTaskHandler).Methods("GET")

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(":8080", r))
}
