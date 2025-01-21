package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pet1/internal/db"
	"pet1/internal/handlers"
	"pet1/internal/taskService"
)

func main() {
	// Инициализация базы данных
	db.InitDB()
	if err := db.DB.AutoMigrate(&taskService.Task{}); err != nil {
		log.Fatal(err)
	}

	repo := taskService.NewTaskRepository(db.DB)
	service := taskService.NewService(repo)

	handler := handlers.NewHandler(service)
	// Маршруты
	r := mux.NewRouter()
	r.HandleFunc("/tasks", handler.PostTaskHandler).Methods("POST")
	r.HandleFunc("/tasks", handler.GetTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", handler.UpdateTaskHandler).Methods("PATCH")
	r.HandleFunc("/tasks/{id}", handler.DeleteTaskHandler).Methods("DELETE")

	// Запускаем сервер
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
