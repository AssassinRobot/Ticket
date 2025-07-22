package rest

import (
	"user/internal/ports"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	api ports.APIPort
}

func NewUserHandler(api ports.APIPort)*UserHandler{
	return &UserHandler{
		api: api,
	}
}


type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}


func (h *UserHandler)ListUsers(ctx *fiber.Ctx) error {
	users, err := h.api.ListUsers(ctx.Context())

	if err != nil {		
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list users: " + err.Error(),
		})
	}

	return ctx.JSON(users)
}

func (h *UserHandler)GetUserByID(ctx *fiber.Ctx) error{
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID: " + err.Error(),
		})
	}

	user, err := h.api.GetUserByID(ctx.Context(), uint(ID))

	if err != nil {		
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user: " + err.Error(),
		})
	}

	return ctx.JSON(user)
}

func (h *UserHandler)SaveUser(ctx *fiber.Ctx) error {
	var createUserRequest = &CreateUserRequest{}
	err := ctx.BodyParser(createUserRequest)
	
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	user, err := h.api.Register(ctx.Context(), createUserRequest.FirstName, createUserRequest.LastName, createUserRequest.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler)UpdateUser(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID: " + err.Error(),
		})
	}

	var UpdateUserRequest = &UpdateUserRequest{}
	err = ctx.BodyParser(UpdateUserRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.api.UpdateUser(ctx.Context(), uint(ID), UpdateUserRequest.FirstName, UpdateUserRequest.LastName)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusNoContent).JSON(nil)
}


func (h *UserHandler)DeleteUser(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID: " + err.Error(),
		})
	}

	err = h.api.DeleteUser(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusNoContent).JSON(nil)
}