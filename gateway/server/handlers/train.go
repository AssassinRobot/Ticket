package handlers

import (
	"gateway/events"

	"github.com/gofiber/fiber/v2"
)

type TrainHandler struct {
	requestHandler events.RequestSender
}

func NewTrainHandler(requestSender events.RequestSender) *TrainHandler {
	return &TrainHandler{
		requestHandler: requestSender,
	}
}

type Train struct {
	Name     string `json:"name"`
	Capacity uint   `json:"capacity"`
}

type CreateTrainRequest struct {
	Train
}

type UpdateTrainRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TrainTravelDetailsRequest struct {
	TrainID       uint   `json:"train_id"`
	Destination   string `json:"destination"`
	Origin        string `json:"origin"`
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
}

func (h *TrainHandler) ListTrains(ctx *fiber.Ctx) error {
	filters := ctx.Queries()

	if len(filters) == 0 {
		trains, err := h.requestHandler.ListTrains(ctx.Context())

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to list trains: " + err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(trains)
	}

	trains, err := h.requestHandler.ListTrainsFiltered(ctx.Context(), filters)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list filtered trains: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(trains)
}

func (h *TrainHandler) ListTrainsSeats(ctx *fiber.Ctx) error {
	trainID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid train ID: " + err.Error(),
		})
	}

	seats, err := h.requestHandler.ListSeatsByTrainID(ctx.Context(), uint(trainID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list seats: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(seats)
}

func (h *TrainHandler) GetTrainByID(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid train ID: " + err.Error(),
		})
	}

	train, err := h.requestHandler.GetTrainByID(ctx.Context(), uint(ID))

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get train: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(train)
}

func (h *TrainHandler) ListTrainTickets(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid train ID: " + err.Error(),
		})
	}

	tickets, err := h.requestHandler.ListTicketsByTrainID(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user tickets: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(tickets)
}

func (h *TrainHandler) CreateTrain(ctx *fiber.Ctx) error {
	var createTrainRequest = &CreateTrainRequest{}
	err := ctx.BodyParser(createTrainRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.requestHandler.CreateTrain(ctx.Context(), createTrainRequest.Name, createTrainRequest.Capacity)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create train: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Train created successfully",
	})
}

func (h *TrainHandler) UpdateTrain(ctx *fiber.Ctx) error {
	var updateTrainRequest = &UpdateTrainRequest{}
	err := ctx.BodyParser(updateTrainRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.requestHandler.UpdateTrainName(ctx.Context(), updateTrainRequest.ID, updateTrainRequest.Name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update train: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "Train updated successfully",
		},
	)
}

func (h *TrainHandler) UpdateTrainTravelDetails(ctx *fiber.Ctx) error {
	var travelDetailsRequest = &TrainTravelDetailsRequest{}
	err := ctx.BodyParser(travelDetailsRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.requestHandler.UpdateTrainTravelDetails(ctx.Context(), travelDetailsRequest.TrainID, travelDetailsRequest.Destination, travelDetailsRequest.Origin, travelDetailsRequest.DepartureTime, travelDetailsRequest.ArrivalTime)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update train travel details: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "Train travel details updated successfully",
		},
	)
}

func (h *TrainHandler) DeleteTrain(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid train ID: " + err.Error(),
		})
	}

	err = h.requestHandler.DeleteTrainByID(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete train: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "Train deleted successfully",
		},
	)
}
