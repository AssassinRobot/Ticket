package rest

import (
	"user/internal/ports"

	"github.com/gofiber/fiber/v2"
)

func InitRouters(api ports.APIPort) *fiber.App {
	userHandler := NewUserHandler(api)

	r := fiber.New()

	v1 := r.Group("/api/v1/users")
	v1.Get("", userHandler.ListUsers)
	v1.Get("/:id", userHandler.GetUserByID)
	v1.Post("", userHandler.SaveUser)
	v1.Patch("/:id", userHandler.UpdateUser)
	v1.Delete("/:id", userHandler.DeleteUser)

	return r
}

