package postgres

import (
	"context"
	"fmt"
	"strings"
	"train/internal/application/core/domain"

	"gorm.io/gorm"
)

type PostgresDBAdapter struct {
	db *gorm.DB
}

func NewPostgresDBAdapter(db *gorm.DB) *PostgresDBAdapter {
	return &PostgresDBAdapter{
		db: db,
	}
}

func (r *PostgresDBAdapter) CreateTrain(ctx context.Context, train *domain.Train) error {
	train.AvailableSeats = train.Capacity

	if err := r.db.WithContext(ctx).Create(train).Error; err != nil {
		return fmt.Errorf("failed to create train: %w", err)
	}

	return nil
}

func (r *PostgresDBAdapter) GetTrainByID(ctx context.Context, ID uint) (*domain.Train, error) {
	var train domain.Train
	err := r.db.WithContext(ctx).Preload("Seats", "booked = false").First(&train, ID).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get train by ID: %w", err)
	}
	return &train, nil
}

func (r *PostgresDBAdapter) ListTrains(ctx context.Context) ([]domain.Train, error) {
	var trains []domain.Train

	err := r.db.WithContext(ctx).Find(&trains).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list trains: %w", err)
	}
	return trains, nil
}

func (r *PostgresDBAdapter) ListTrainsFiltered(ctx context.Context, filters *domain.TrainFilters) ([]domain.Train, error) {
	var trains []domain.Train

	query := r.db.WithContext(ctx).Model(&domain.Train{})

	var conditions []string
	var args []interface{}

	if filters.Name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+filters.Name+"%")
	}

	if filters.Origin != "" {
		conditions = append(conditions, "destination = ?")
		args = append(args, filters.Destination)
	}

	if filters.Destination != "" {
		conditions = append(conditions, "destination = ?")
		args = append(args, filters.Destination)
	}

	if filters.AvailableSeats != 0 {
		conditions = append(conditions, "available_seats >= ?")
		args = append(args, filters.AvailableSeats)
	}

	if !filters.DepartureTime.IsZero() {
		conditions = append(conditions, "departure_time = ?")
		args = append(args, filters.DepartureTime)
	}

	if !filters.ArrivalTime.IsZero() {
		conditions = append(conditions, "arrival_time = ?")
		args = append(args, filters.ArrivalTime)
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	query = query.Preload("Seats", "booked = false")

	err := query.Find(&trains).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no train available")

		}
		return nil, fmt.Errorf("search trains: args:  %v\n: %w", args, err)

	}

	return trains, nil
}

func (r *PostgresDBAdapter) IsTrainAvailable(ctx context.Context, ID uint) (bool, error) {
	var isFull bool

	err := r.db.WithContext(ctx).Model(&domain.Train{}).Where("id = ?", ID).Select("is_full").Scan(&isFull).Error
	if err != nil {
		return false, fmt.Errorf("failed get train column: %w", err)
	}

	return isFull, nil
}

func (r *PostgresDBAdapter) UpdateTrain(ctx context.Context, ID uint, name string) error {
	err := r.db.WithContext(ctx).Model(&domain.Train{}).Where("id = ?", ID).Update("name", name).Error
	if err != nil {
		return fmt.Errorf("failed to update train name: train ID:%d %w", ID, err)
	}
	
	return nil
}

func (r *PostgresDBAdapter) UpdateTrainTravelDetails(ctx context.Context, ID uint, travelDetails *domain.TrainTravelDetails) error {
	train, err := r.GetTrainByID(ctx, ID)
	if err != nil {
		return err
	}

	train.Destination = travelDetails.Destination
	train.Origin = travelDetails.Origin
	train.DepartureTime = travelDetails.DepartureTime
	train.ArrivalTime = travelDetails.ArrivalTime

	err = r.db.WithContext(ctx).Save(&train).Error
	if err != nil {
		return fmt.Errorf("failed to update train travel details: %w", err)
	}

	return nil
}

func (r *PostgresDBAdapter) MinusTrainAvailableSeats(ctx context.Context, ID uint) error {
	train, err := r.GetTrainByID(ctx, ID)
	if err != nil {
		return err
	}

	if train.AvailableSeats <= 0 {
		return fmt.Errorf("no available seats to book")
	}

	train.AvailableSeats--
	if train.AvailableSeats == 0 {
		train.IsFull = true
	}

	err = r.db.WithContext(ctx).Save(&train).Error

	if err != nil {
		return fmt.Errorf("failed to update train available seats: %w", err)
	}

	return nil
}

func (r *PostgresDBAdapter) PlusTrainAvailableSeats(ctx context.Context, ID uint) error {
	train, err := r.GetTrainByID(ctx, ID)
	if err != nil {
		return err
	}

	if train.AvailableSeats == train.Capacity {
		return fmt.Errorf("all seats are already available")
	}

	train.AvailableSeats++
	if train.AvailableSeats > 0 {
		train.IsFull = false
	}

	err = r.db.WithContext(ctx).Save(&train).Error

	if err != nil {
		return fmt.Errorf("failed to update train available seats: %w", err)
	}

	return nil
}

func (r *PostgresDBAdapter) DeleteTrain(ctx context.Context, ID uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Train{}, "id = ?", ID)
	if result.Error != nil {
		return fmt.Errorf("delete train by id: train ID:%d \n%w", ID, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("train not found with ID: %d", ID)
	}

	return nil
}

func (r *PostgresDBAdapter) CreateSeat(ctx context.Context, seat *domain.Seat) error {
	err := r.db.WithContext(ctx).Create(seat).Error
	if err != nil {
		return fmt.Errorf("failed to create seat: %w", err)
	}

	return nil
}

func (r *PostgresDBAdapter) UpdateSeatNumber(ctx context.Context, seatID uint, seatNumber uint) error {
	err := r.db.WithContext(ctx).Model(&domain.Seat{}).Where("id = ?", seatID).Update("seat_number", seatNumber).Error
	if err != nil {
		return fmt.Errorf("failed to update seat number: seat ID:%d %w", seatID, err)
	}
	return nil
}


func (r *PostgresDBAdapter) UpdateSeatUser(ctx context.Context, seatID,userID uint) error {
	err := r.db.WithContext(ctx).Model(&domain.Seat{}).Where("id = ?", seatID).Update("user_id", userID).Error
	if err != nil {
		return fmt.Errorf("failed to update user ID: seat ID:%d %w", seatID, err)
	}
	return nil
}

func (r *PostgresDBAdapter) GetSeatByID(ctx context.Context, id uint) (*domain.Seat, error) {
	seat := &domain.Seat{}

	err := r.db.WithContext(ctx).First(seat, id).Error
	if err != nil {
		return nil, fmt.Errorf("get seat by id: seat ID: %d \n%w", id, err)
	}

	return seat, nil
}

func (r *PostgresDBAdapter) ListSeatsByTrainID(ctx context.Context, trainID uint) ([]domain.Seat, error) {
	var seats []domain.Seat

	err := r.db.WithContext(ctx).Where("train_id = ?", trainID).Find(&seats).Error
	if err != nil {
		return nil, fmt.Errorf("list seats by train ID: %w", err)
	}

	return seats, nil
}

func (r *PostgresDBAdapter) UpdateSeatBookingStatus(ctx context.Context, seatID uint, booked bool) error {
	err := r.db.WithContext(ctx).Model(&domain.Seat{}).Where("id = ?", seatID).Update("booked", booked).Error

	if err != nil {
		return fmt.Errorf("failed to update booking status: seat ID:%d %w", seatID, err)
	}

	return nil
}

func (r *PostgresDBAdapter) DeleteSeat(ctx context.Context, ID uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Seat{}, ID)
	if result.Error != nil {
		return fmt.Errorf("delete seat by id: seat ID:%d \n%w", ID, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("seat not found with ID: %d", ID)
	}

	return nil
}
