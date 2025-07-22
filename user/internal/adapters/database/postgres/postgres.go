package postgres

import (
	"context"
	"errors"
	"fmt"
	"user/internal/application/core/domain"
	"user/utils"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type DatabasePostgresAdapter struct {
	db *gorm.DB
}

func NewDatabasePostgresAdapter(db *gorm.DB) *DatabasePostgresAdapter {
	return &DatabasePostgresAdapter{
		db: db,
	}
}

func (u *DatabasePostgresAdapter) SaveUser(ctx context.Context, user *domain.User) error {
	err := u.db.Create(user).Error
	if err != nil {
		if utils.CheckErrorForWord(err, "email") {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("save user: %v \n%w", user, err)
	}

	return nil

}

func (u *DatabasePostgresAdapter) GetUserByID(ctx context.Context, ID uint) (*domain.User, error) {
	user := &domain.User{}

	err := u.db.WithContext(ctx).First(user, ID).Error
	if err != nil {
		return nil, fmt.Errorf("get user by id: user ID:  %d \n%w", ID, err)
	}

	return user, nil
}

func (u *DatabasePostgresAdapter) ListUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	err := u.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, nil
}

func (u *DatabasePostgresAdapter) UpdateUser(ctx context.Context, ID uint, firstName, lastName string) (*domain.User,error) {
	user, err := u.GetUserByID(ctx, ID)
	if err != nil {
		return nil,err
	}

	user.FirstName = firstName
	user.LastName = lastName

	err = u.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return  nil,fmt.Errorf("update user by id: user ID:%d ,first name:%s ,last name:%s \n%w", ID, firstName, lastName, err)
	}

	return user,nil
}

func (u *DatabasePostgresAdapter) DeleteUser(ctx context.Context, ID uint) error {

	result := u.db.WithContext(ctx).Delete(&domain.User{}, ID)
	if result.Error != nil {
		return fmt.Errorf("delete user by id: user ID:%d \n%w", ID, result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
