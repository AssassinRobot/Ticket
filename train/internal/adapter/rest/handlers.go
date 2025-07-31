package rest

import (
	"train/internal/ports"

	"github.com/gofiber/fiber/v2"
)

type TrainHandler struct {
	api ports.APIPort
}

func NewTrainHandler(api ports.APIPort) *TrainHandler {
	return &TrainHandler{
		api: api,
	}
}

type Train struct {
	Name     string `json:"name"`
	Capacity uint   `json:"capacity"`
}

type CreateTrainRequest struct {
	Train
}

type CreateSeatRequest struct {
	TrainID    uint `json:"train_id"`
	SeatNumber uint `json:"seat_number"`
}

type UpdateTrainRequest struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}

type UpdateSeatRequest struct {
	ID         uint `json:"id"`
	SeatNumber uint `json:"seat_number"`
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
		trains, err := h.api.ListTrains(ctx.Context())

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to list trains: " + err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(trains)
	}

	trains, err := h.api.ListTrainsFiltered(ctx.Context(), filters)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list filtered trains: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(trains)
}

func (h *TrainHandler) GetTrainByID(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid train ID: " + err.Error(),
		})
	}

	train, err := h.api.GetTrainByID(ctx.Context(), uint(ID))

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get train: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(train)
}

func (h *TrainHandler) CreateTrain(ctx *fiber.Ctx) error {
	var createTrainRequest = &CreateTrainRequest{}
	err := ctx.BodyParser(createTrainRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.api.CreateTrain(ctx.Context(), createTrainRequest.Name, createTrainRequest.Capacity)
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

	err = h.api.UpdateTrain(ctx.Context(), updateTrainRequest.ID, updateTrainRequest.Name)
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

	err = h.api.UpdateTrainTravelDetails(ctx.Context(), travelDetailsRequest.TrainID, travelDetailsRequest.Destination, travelDetailsRequest.Origin, travelDetailsRequest.DepartureTime, travelDetailsRequest.ArrivalTime)
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

	err = h.api.DeleteTrain(ctx.Context(), uint(ID))
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

func (h *TrainHandler) CreateSeat(ctx *fiber.Ctx) error {
	var createSeatRequest = &CreateSeatRequest{}
	err := ctx.BodyParser(createSeatRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.api.CreateSeat(ctx.Context(), createSeatRequest.TrainID, createSeatRequest.SeatNumber)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create seat: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Seat created successfully",
	})
}

func (h *TrainHandler) UpdateSeatNumber(ctx *fiber.Ctx) error {
	var updateSeatRequest = &UpdateSeatRequest{}
	err := ctx.BodyParser(updateSeatRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.api.UpdateSeatNumber(ctx.Context(), updateSeatRequest.ID, updateSeatRequest.SeatNumber)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update seat number: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Seat number updated successfully",
	})
}

func (h *TrainHandler) GetSeatByID(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid seat ID: " + err.Error(),
		})
	}

	seat, err := h.api.GetSeatByID(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get seat: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(seat)
}

func (h *TrainHandler) ListSeatsByTrainID(ctx *fiber.Ctx) error {
	trainID, err := ctx.ParamsInt("train_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid train ID: " + err.Error(),
		})
	}

	seats, err := h.api.ListSeatsByTrainID(ctx.Context(), uint(trainID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list seats: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(seats)
}

func (h *TrainHandler) DeleteSeat(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Seat ID: " + err.Error(),
		})
	}

	err = h.api.DeleteSeat(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete Seat: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "Seat deleted successfully",
		},
	)
}
