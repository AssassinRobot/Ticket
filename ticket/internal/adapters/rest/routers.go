package rest

import (
	"ticket/internal/ports"

	"github.com/gofiber/fiber/v2"
)

func InitRouters(api ports.APIPort) *fiber.App {

	r := fiber.New()

	v1 := r.Group("/api/v1")

	ticketHandler := NewTicketHandler(api)
	trainRoutes := v1.Group("/tickets")
	trainRoutes.Get("/:id", ticketHandler.GetTicketByID)
	trainRoutes.Get("/user/:user_id", ticketHandler.GetTicketsByUserID)
	trainRoutes.Post("/book", ticketHandler.BookTicket)
	trainRoutes.Post("/cancel", ticketHandler.CancelTicket)

	return r
}
