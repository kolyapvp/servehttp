package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
	UserID uint   `json:"user_id"`
}

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Tasks    []Task `json:"tasks" gorm:"foreignKey:UserID"`
}
