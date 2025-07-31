package domain

import (
	"time"
)

type Ticket struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	TrainID    uint
	SeatID uint
	SeatNumber uint
	TravelDetails
	ExpiresAt  time.Time
	CanceledAt *time.Time 
	CreatedAt  time.Time
}

func NewTicket(userID, trainID,seatID , seatNumber uint, expiresAt time.Time, travelDetails TravelDetails) *Ticket {
	return &Ticket{
		UserID:        userID,
		TrainID:       trainID,
		SeatID: seatID,
		SeatNumber:    seatNumber,
		TravelDetails: travelDetails,
		ExpiresAt:     expiresAt,
	}
}
