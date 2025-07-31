package domain

import "time"

type Train struct {
	ID uint
	Name string
    Seats []Seat
	TravelDetails
}

type TravelDetails struct {
	Destination   string
	Origin        string
	DepartureTime time.Time
	ArrivalTime   time.Time
}

type Seat struct{
	ID uint
	TrainID uint
	UserID uint
    SeatNumber uint
}


