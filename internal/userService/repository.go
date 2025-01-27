package userService

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	GetUserByID(id uint) (User, error)
	UpdateUserByID(id uint, user User) (User, error)
	DeleteUserByID(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user User) (User, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) GetUserByID(id uint) (User, error) {
	var user User
	result := r.db.Preload("Tasks").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, errors.New("user not found")
		}
		return User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) UpdateUserByID(id uint, user User) (User, error) {
	var existingUser User
	result := r.db.First(&existingUser, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, errors.New("user not found")
		}
		return User{}, result.Error
	}

	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Password != "" {
		existingUser.Password = user.Password
	}

	saveResult := r.db.Save(&existingUser)
	if saveResult.Error != nil {
		return User{}, saveResult.Error
	}

	return existingUser, nil
}

func (r *userRepository) DeleteUserByID(id uint) error {
	var user User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return result.Error
	}

	deleteResult := r.db.Delete(&user)
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	if deleteResult.RowsAffected == 0 {
		return errors.New("no user was deleted")
	}

	return nil
}
