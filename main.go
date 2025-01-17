package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// Глобальная переменная для хранения task
var task string

// POST handler для обновления task
func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON из тела запроса
	var requestBody struct {
		Task string `json:"task"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Обновляем глобальную переменную task
	task = requestBody.Task

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Task updated successfully")
}

// GET handler для возвращения приветствия с task
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем текущую task
	currentTask := task

	// Возвращаем приветствие
	fmt.Fprintf(w, "hello, %s\n", currentTask)
}

func main() {

	router := mux.NewRouter()
	// Регистрируем обработчики
	router.HandleFunc("/get", getTaskHandler).Methods("GET")
	router.HandleFunc("/post", postTaskHandler).Methods("POST")

	http.ListenAndServe(":8080", router)
}
