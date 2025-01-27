package handlers

import (
	"context"
	"fmt"
	"pet1/internal/userService"
	"pet1/internal/web/users"
)

// UserHandler структура для обработки запросов пользователей
type UserHandler struct {
	Service *userService.UserService
}

func NewUserHandler(service *userService.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

// GetUsers реализует получение всех пользователей
func (h *UserHandler) GetUsers(_ context.Context, _ users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	allUsers, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	response := users.GetUsers200JSONResponse{}
	for _, usr := range allUsers {
		user := users.User{
			Id:        &usr.ID,
			Email:     &usr.Email,
			Password:  &usr.Password,
			CreatedAt: &usr.CreatedAt,
			UpdatedAt: &usr.UpdatedAt,
		}
		response = append(response, user)
	}

	return response, nil
}

// PostUsers реализует создание нового пользователя
func (h *UserHandler) PostUsers(_ context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	userRequest := request.Body
	userToCreate := userService.User{
		Email:    *userRequest.Email,
		Password: *userRequest.Password,
	}

	createdUser, err := h.Service.CreateUser(userToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := users.PostUsers201JSONResponse{
		Id:        &createdUser.ID,
		Email:     &createdUser.Email,
		Password:  &createdUser.Password,
		CreatedAt: &createdUser.CreatedAt,
		UpdatedAt: &createdUser.UpdatedAt,
	}

	return response, nil
}

// DeleteUsersId реализует удаление пользователя по ID
func (h *UserHandler) DeleteUsersId(_ context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
	id := request.Id

	err := h.Service.DeleteUserByID(id)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "no user was deleted" {
			return users.DeleteUsersId404Response{}, nil
		}
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return users.DeleteUsersId204Response{}, nil
}

// PatchUsersId реализует обновление пользователя по ID
func (h *UserHandler) PatchUsersId(_ context.Context, request users.PatchUsersIdRequestObject) (users.PatchUsersIdResponseObject, error) {
	id := request.Id
	body := request.Body

	userToUpdate := userService.User{}

	if body.Email != nil {
		userToUpdate.Email = *body.Email
	}
	if body.Password != nil {
		userToUpdate.Password = *body.Password
	}

	updatedUser, err := h.Service.UpdateUserByID(id, userToUpdate)
	if err != nil {
		if err.Error() == "user not found" {
			return users.PatchUsersId404Response{}, nil
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := users.PatchUsersId200JSONResponse{
		Id:        &updatedUser.ID,
		Email:     &updatedUser.Email,
		Password:  &updatedUser.Password,
		CreatedAt: &updatedUser.CreatedAt,
		UpdatedAt: &updatedUser.UpdatedAt,
	}

	return response, nil
}
