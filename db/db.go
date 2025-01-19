package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// переменная, через которую мы будем работать с БД
var DB *gorm.DB

func InitDB() {
	// в dsn вводим данные, которые мы указали при создании контейнера
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Выполняем SQL-запрос для создания таблицы tasks, если она не существует
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		task TEXT NOT NULL,
		is_done BOOLEAN NOT NULL DEFAULT FALSE
	);
	`
	if err := DB.Exec(createTableQuery).Error; err != nil {
		log.Fatal("Failed to create tasks table: ", err)
	}

	log.Println("Database initialized successfully")
}
