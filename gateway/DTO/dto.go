package dto

import "time"

// DTO types
type (
	TicketDTO struct {
		ID         uint
		UserID     uint
		TrainID    uint
		SeatNumber uint
		TravelDetailDTO
		ExpiresAt time.Time
		CanceledAt *time.Time
	}

	TrainDTO struct {
		ID       uint
		Name     string
		Capacity uint
		IsFull   bool
		Seats    []SeatDTO
		TravelDetailDTO
	}

	TravelDetailDTO struct {
		Origin        string
		Destination   string
		DepartureTime time.Time
		ArrivalTime   time.Time
	}

	SeatDTO struct {
		ID      uint
		Number  uint
		TrainID uint
		UserID  uint
	}

	UserDTO struct {
		ID        uint
		FirstName string
		LastName  string
		Email     string
	}
)
