package rest

import (
	"ticket/internal/ports"

	"github.com/gofiber/fiber/v2"
)

type TicketHandler struct {
	api ports.APIPort
}

func NewTicketHandler(api ports.APIPort) *TicketHandler {
	return &TicketHandler{
		api: api,
	}
}

type BookTicketRequest struct {
	UserID     uint `json:"user_id"`
	TrainID uint   `json:"train_id"`
	TicketNumber uint `json:"ticket_number"`
}

type CancelTicketRequest struct {
	TicketID uint `json:"ticket_id"`
}



func (h *TicketHandler) GetTicketByID(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID: " + err.Error(),
		})
	}

	ticket, err := h.api.GetTicketByID(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get ticket" + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(ticket)
}

func (h *TicketHandler) GetTicketsByUserID(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID: " + err.Error(),
		})
	}

	tickets, err := h.api.GetTicketsByUserID(ctx.Context(), uint(ID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user tickets: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(tickets)
}

func (h *TicketHandler) BookTicket(ctx *fiber.Ctx) error {
	var bookTicketRequest = &BookTicketRequest{}
	err := ctx.BodyParser(bookTicketRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	tickets,err := h.api.BookTicket(ctx.Context(), bookTicketRequest.UserID, bookTicketRequest.TrainID,bookTicketRequest.TicketNumber)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to book ticket: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "tickets booked successfully",
		"tickets":tickets,
	})
}

func (h *TicketHandler) CancelTicket(ctx *fiber.Ctx) error {
	var cancelTicketRequest = &CancelTicketRequest{}
	err := ctx.BodyParser(cancelTicketRequest)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	err = h.api.CancelTicket(ctx.Context(), cancelTicketRequest.TicketID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel ticket: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"message": "ticket canceled successfully",
		},
	)
}
