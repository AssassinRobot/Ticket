package router

import (
	"gateway/events"
	"gateway/server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitRouters(requestSender events.RequestSender) *fiber.App {
	r := fiber.New()

	r.Use(logger.New())

	v1 := r.Group("/api/v1")

	{
		userHandler := handlers.NewUserHandler(requestSender)
		r := v1.Group("/users")

		r.Get("", userHandler.ListUsers)
		r.Get("/:id", userHandler.GetUserByID)
		r.Get("/:id/tickets", userHandler.ListUserTickets)
		r.Post("", userHandler.CreateUser)
		r.Patch("", userHandler.UpdateUser)
		r.Delete("/:id", userHandler.DeleteUser)
	}

	{
		trainHandler := handlers.NewTrainHandler(requestSender)

		r := v1.Group("/trains")

		r.Get("/", trainHandler.ListTrains)
		r.Get("/:id", trainHandler.GetTrainByID)
		r.Get("/:id/seats", trainHandler.ListTrainsSeats)
		r.Get("/:id/tickets", trainHandler.ListTrainTickets)
		r.Post("", trainHandler.CreateTrain)
		r.Patch("", trainHandler.UpdateTrain)
		r.Patch("/travel-details", trainHandler.UpdateTrainTravelDetails)
		r.Delete("/:id", trainHandler.DeleteTrain)
	}

	{
		seatHandler := handlers.NewSeatHandler(requestSender)

		r := v1.Group("/seats")

		r.Get("/:id", seatHandler.GetSeatByID)
		r.Post("", seatHandler.CreateSeat)
		r.Patch("", seatHandler.UpdateSeatNumber)
		r.Delete("/:id", seatHandler.DeleteSeat)
	}

	{
		ticketHandler := handlers.NewTicketHandler(requestSender)

		r := v1.Group("/tickets")

		r.Get("/:id", ticketHandler.GetTicketByID)
		r.Post("/book", ticketHandler.BookTicket)
		r.Post("/cancel", ticketHandler.CancelTicket)

	}

	return r

}
