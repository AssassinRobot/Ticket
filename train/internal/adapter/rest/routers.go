package rest

import (
	"train/internal/ports"

	"github.com/gofiber/fiber/v2"
)

func InitRouters(api ports.APIPort) *fiber.App {

	r := fiber.New()

	v1 := r.Group("/api/v1")

	trainHandler := NewTrainHandler(api)
	trainRoutes := v1.Group("/trains")
	trainRoutes.Get("/", trainHandler.ListTrains)
	trainRoutes.Get("/:id", trainHandler.GetTrainByID)
	trainRoutes.Post("", trainHandler.CreateTrain)
	trainRoutes.Patch("/", trainHandler.UpdateTrain)
	trainRoutes.Patch("/travel-details", trainHandler.UpdateTrainTravelDetails)
	trainRoutes.Delete("/:id", trainHandler.DeleteTrain)

	seatRoutes := v1.Group("/seats")
	seatRoutes.Get("/:id", trainHandler.GetSeatByID)
	seatRoutes.Get("/train/:train_id", trainHandler.ListSeatsByTrainID)
	seatRoutes.Post("", trainHandler.CreateSeat)
	seatRoutes.Patch("", trainHandler.UpdateSeatNumber)
	seatRoutes.Delete("/:id", trainHandler.DeleteSeat)

	return r
}
