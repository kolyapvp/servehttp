package userService

import (
	"pet1/internal/taskService"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Tasks    []taskService.Task `json:"tasks" gorm:"foreignKey:UserID"`
}

type Task struct {
	gorm.Model
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
	UserID uint   `json:"user_id"`
}
