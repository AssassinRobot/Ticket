package handlers

import (
	"gateway/events"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	requestHandler events.RequestSender
}

func NewUserHandler(requestHandler events.RequestSender) *UserHandler {
	return &UserHandler{
		requestHandler: requestHandler,
	}
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateUserRequest struct {
	ID uint `json:"ID"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *UserHandler) ListUsers(ctx *fiber.Ctx) error {
	users, err := h.requestHandler.ListUsers(ctx.Context())

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list users: " + err.Error(),
		})
	}

	return ctx.JSON(users)
}

func (h *UserHandler) GetUserByID(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID: " + err.Error(),
		})
	}

	user, err := h.requestHandler.GetUserByID(ctx.Context(), uint(ID))

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user: " + err.Error(),
		})
	}

	return ctx.JSON(user)
}


func (h *UserHandler) CreateUser(ctx *fiber.Ctx) error {
	var createUserRequest = &CreateUserRequest{}
	err := ctx.BodyParser(createUserRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.requestHandler.CreateUser(ctx.Context(), createUserRequest.FirstName, createUserRequest.LastName, createUserRequest.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg": "user created successfully",
	})
}

func (h *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
	var updateUserRequest = &UpdateUserRequest{}
	err := ctx.BodyParser(updateUserRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.requestHandler.UpdateUserByID(ctx.Context(), uint(updateUserRequest.ID), updateUserRequest.FirstName, updateUserRequest.LastName)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":"user was updated successfully",
	})
}

func (h *UserHandler) DeleteUser(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID: " + err.Error(),
		})
	}

	err = h.requestHandler.DeleteUserByID(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":"user was deleted successfully",
	})
}
