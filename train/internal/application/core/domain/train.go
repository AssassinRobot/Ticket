package domain

import "time"

type Train struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	TrainTravelDetails
	Capacity       uint32 // Total number of seats in the train
	AvailableSeats uint32 // Number of seats available for booking
	IsFull         bool   // Indicates if the train is fully booked
	Seats          []Seat `gorm:"foreignKey:TrainID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TrainFilters struct {
	Name           string
	AvailableSeats uint
	TrainTravelDetails
}

type TrainTravelDetails struct {
	Destination   string
	Origin        string
	DepartureTime time.Time
	ArrivalTime   time.Time
}

func NewTrain(name string, capacity uint32) *Train {
	return &Train{
		Name:     name,
		Capacity: capacity,
	}
}

func NewFilterTrains(name, origin, destination string, departureTime, arrivalTime time.Time) *TrainFilters {
	return &TrainFilters{
		Name:          name,
		TrainTravelDetails: TrainTravelDetails{
			Origin:        origin,
			Destination:   destination,
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
		},
	}
}
