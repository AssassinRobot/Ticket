package domain

import "time"

type Seat struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint // Foreign key to User
	TrainID    uint // Foreign key to Train
	SeatNumber uint
	Booked     bool // Indicates if the seat is booked
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewSeat(trainID uint, seatNumber uint) *Seat {
	return &Seat{
		TrainID:    trainID,
		SeatNumber: seatNumber,
	}
}
