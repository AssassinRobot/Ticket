package handlers

import (
	"gateway/events"

	"github.com/gofiber/fiber/v2"
)

type SeatHandler struct {
	requestHandler events.RequestSender
}

type CreateSeatRequest struct {
	TrainID    uint `json:"train_id"`
	SeatNumber uint `json:"seat_number"`
}

type UpdateSeatRequest struct {
	ID         uint `json:"id"`
	SeatNumber uint `json:"seat_number"`
}

func NewSeatHandler(requestSender events.RequestSender) *SeatHandler {
	return &SeatHandler{
		requestHandler: requestSender,
	}
}

func (h *SeatHandler) CreateSeat(ctx *fiber.Ctx) error {
	var createSeatRequest = &CreateSeatRequest{}
	err := ctx.BodyParser(createSeatRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.requestHandler.CreateSeat(ctx.Context(), createSeatRequest.TrainID, createSeatRequest.SeatNumber)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create seat: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Seat created successfully",
	})
}

func (h *SeatHandler) UpdateSeatNumber(ctx *fiber.Ctx) error {
	var updateSeatRequest = &UpdateSeatRequest{}
	err := ctx.BodyParser(updateSeatRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.requestHandler.UpdateSeatNumberBySeatID(ctx.Context(), updateSeatRequest.ID, updateSeatRequest.SeatNumber)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update seat number: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Seat number updated successfully",
	})
}

func (h *SeatHandler) GetSeatByID(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid seat ID: " + err.Error(),
		})
	}

	seat, err := h.requestHandler.GetSeatByID(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get seat: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(seat)
}



func (h *SeatHandler) DeleteSeat(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Seat ID: " + err.Error(),
		})
	}

	err = h.requestHandler.DeleteSeatBySeatID(ctx.Context(), uint(ID))
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
