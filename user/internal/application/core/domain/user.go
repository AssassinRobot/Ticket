package domain

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(firstName, lastName, email  string) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
}