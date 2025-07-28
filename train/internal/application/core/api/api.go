package api

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"train/internal/application/core/domain"
	"train/internal/ports"
)

type APIAdapter struct {
	DatabasePort ports.DatabasePort
}

func NewAPIAdapter(DBport ports.DatabasePort)*APIAdapter{
	return  &APIAdapter{
		DatabasePort:DBport,
	}
}

func (a *APIAdapter) CreateTrain(ctx context.Context, name string, capacity uint) error {
	train := domain.NewTrain(name, uint32(capacity))

	return a.DatabasePort.CreateTrain(ctx, train)
}

func (a *APIAdapter) GetTrainByID(ctx context.Context, ID uint) (*domain.Train, error) {
	return a.DatabasePort.GetTrainByID(ctx, ID)
}

func (a *APIAdapter) ListTrains(ctx context.Context) ([]domain.Train, error) {
	return a.DatabasePort.ListTrains(ctx)
}

func (a *APIAdapter) ListTrainsFiltered(ctx context.Context, filters map[string]string) ([]domain.Train, error) {
	var trainFilters = &domain.TrainFilters{}

	for key, value := range filters {
		switch key {
		case "name":
			trainFilters.Name = value
		case "available_seats":
			availableSeats, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse 'available_seats' filter: %w", err)
			}

			trainFilters.AvailableSeats = uint(availableSeats)
		case "origin":
			trainFilters.Origin = value
		case "destination":
			trainFilters.Destination = value
		case "departure_time":
			timeParsed, err := time.Parse(time.DateTime, value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse 'departure_time' filter: %w", err)
			}

			trainFilters.DepartureTime = timeParsed
		case "arrival_time":
			timeParsed, err := time.Parse(time.DateTime, value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse 'arrival_time' filter: %w", err)
			}

			trainFilters.ArrivalTime = timeParsed
		}
	}

	return a.DatabasePort.ListTrainsFiltered(ctx, trainFilters)
}

func (a *APIAdapter) UpdateTrain(ctx context.Context, ID uint, name string, capacity uint32) error {
	return a.DatabasePort.UpdateTrain(ctx, ID, name, capacity)
}

func (a *APIAdapter) UpdateTrainTravelDetails(ctx context.Context, TrainID uint, destination, origin, departureTime, arrivalTime string) error {

	departureTimeParsed, err := time.Parse(time.DateTime, departureTime)
	if err != nil {
		return fmt.Errorf("failed to parse departure time: %w", err)
	}

	arrivalTimeParsed, err := time.Parse(time.DateTime, arrivalTime)
	if err != nil {
		return fmt.Errorf("failed to parse arrival time: %w", err)
	}

	travelDetails := &domain.TrainTravelDetails{
		Destination:   destination,
		Origin:        origin,
		DepartureTime: departureTimeParsed,
		ArrivalTime:   arrivalTimeParsed,
	}

	return a.DatabasePort.UpdateTrainTravelDetails(ctx, TrainID, travelDetails)
}

func (a *APIAdapter) DeleteTrain(ctx context.Context, ID uint) error {
	return a.DatabasePort.DeleteTrain(ctx, ID)
}

func (a *APIAdapter) CreateSeat(ctx context.Context, trainID, seatNumber uint) error {
	seat := domain.NewSeat(trainID, seatNumber)

	return a.DatabasePort.CreateSeat(ctx, seat)
}

func (a *APIAdapter) GetSeatByID(ctx context.Context, ID uint) (*domain.Seat, error) {
	return a.DatabasePort.GetSeatByID(ctx, ID)
}

func (a *APIAdapter) UpdateSeatNumber(ctx context.Context, seatID uint, seatNumber uint) error {
	err := a.DatabasePort.UpdateSeatNumber(ctx, seatID, seatNumber)
	if err != nil {
		return fmt.Errorf("failed to update seat number: %w", err)
	}

	return nil
}

func (a *APIAdapter) SeatBooked(ctx context.Context, seatID uint, trainID uint) error {
	IsTrainAvailable, err := a.DatabasePort.IsTrainAvailable(ctx, trainID)
	if err != nil {
		return err
	}

	if !IsTrainAvailable {
		return fmt.Errorf("train with ID %d is not available", trainID)
	}

	err = a.DatabasePort.UpdateSeatBookingStatus(ctx, seatID, true)
	if err != nil {
		return err
	}

	return a.DatabasePort.MinusTrainAvailableSeats(ctx, trainID)
}

func (a *APIAdapter) CancelSeatBooking(ctx context.Context, seatID uint, trainID uint) error {
	err := a.DatabasePort.UpdateSeatBookingStatus(ctx, seatID, false)
	if err != nil {
		return err
	}

	return a.DatabasePort.PlusTrainAvailableSeats(ctx, trainID)
}

func (a *APIAdapter) ListSeatsByTrainID(ctx context.Context, trainID uint) ([]domain.Seat, error) {
	return a.DatabasePort.ListSeatsByTrainID(ctx, trainID)
}

func (a *APIAdapter) DeleteSeat(ctx context.Context,ID uint)error{
	return  a.DatabasePort.DeleteSeat(ctx,ID)
}