package taskService

import (
	"errors"

	"gorm.io/gorm"
)

type TaskRepository interface {
	// CreateTask - Передаем в функцию task типа Task из orm.go
	// возвращаем созданный Task и ошибку
	CreateTask(task Task) (Task, error)
	// GetAllTasks - Возвращаем массив из всех задач в БД и ошибку
	GetAllTasks() ([]Task, error)
	// UpdateTaskByID - Передаем id и Task, возвращаем обновленный Task
	// и ошибку
	UpdateTaskByID(id uint, task Task) (Task, error)
	// DeleteTaskByID - Передаем id для удаления, возвращаем только ошибку
	DeleteTaskByID(id uint) error
	GetTasksByUserID(userID uint) ([]Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

// (r *taskRepository) привязывает данную функцию к нашему репозиторию
func (r *taskRepository) CreateTask(task Task) (Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// UpdateTaskByID обновляет задачу по ее ID
func (r *taskRepository) UpdateTaskByID(id uint, task Task) (Task, error) {
	// Ищем задачу в базе данных по ID
	var existingTask Task
	result := r.db.First(&existingTask, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Если задача не найдена, возвращаем ошибку
			return Task{}, errors.New("task not found")
		}
		// Если возникла другая ошибка при поиске, возвращаем её
		return Task{}, result.Error
	}

	// Обновляем поля задачи, если они предоставлены
	if task.Task != "" {
		existingTask.Task = task.Task
	}
	existingTask.IsDone = task.IsDone

	// Сохраняем обновленную задачу в базе данных
	saveResult := r.db.Save(&existingTask)
	if saveResult.Error != nil {
		return Task{}, saveResult.Error
	}

	// Возвращаем обновленную задачу и nil (отсутствие ошибки)
	return existingTask, nil
}

// DeleteTaskByID удаляет задачу по ее ID
func (r *taskRepository) DeleteTaskByID(id uint) error {
	// Ищем задачу в базе данных по ID
	var task Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Если задача не найдена, возвращаем ошибку
			return errors.New("task not found")
		}
		// Если возникла другая ошибка при поиске, возвращаем её
		return result.Error
	}

	// Удаляем задачу из базы данных
	deleteResult := r.db.Delete(&task)
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	// Проверяем, была ли удалена хотя бы одна запись
	if deleteResult.RowsAffected == 0 {
		return errors.New("no task was deleted")
	}

	// Возвращаем nil, указывая на отсутствие ошибки
	return nil
}

// GetTasksByUserID получает все задачи пользователя по его ID
func (r *taskRepository) GetTasksByUserID(userID uint) ([]Task, error) {
	var tasks []Task
	result := r.db.Where("user_id = ?", userID).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}
