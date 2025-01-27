package main

import (
	"log"
	"pet1/internal/db"
	"pet1/internal/handlers"
	"pet1/internal/taskService"
	"pet1/internal/userService"
	"pet1/internal/web/tasks"
	"pet1/internal/web/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Инициализация БД
	db.InitDB()

	// Инициализация сервисов задач
	tasksRepo := taskService.NewTaskRepository(db.DB)
	tasksService := taskService.NewService(tasksRepo)
	tasksHandler := handlers.NewTaskHandler(tasksService)

	// Инициализация сервисов пользователей
	usersRepo := userService.NewUserRepository(db.DB)
	usersService := userService.NewService(usersRepo)
	usersHandler := handlers.NewUserHandler(usersService)

	// Инициализируем echo
	e := echo.New()

	// используем Logger и Recover
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Регистрация обработчиков задач
	tasksStrictHandler := tasks.NewStrictHandler(tasksHandler, nil)
	tasks.RegisterHandlers(e, tasksStrictHandler)

	// Регистрация обработчиков пользователей
	usersStrictHandler := users.NewStrictHandler(usersHandler, nil)
	users.RegisterHandlers(e, usersStrictHandler)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}
