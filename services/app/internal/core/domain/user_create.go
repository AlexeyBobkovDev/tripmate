package domain

import "time"

type UserCreate struct {
	Name         string
	Surname      string
	Username     string
	BirthDate    time.Time
	Description  string
	Email        string
	PhoneNumber  string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func NewUserCreate(
	name string,
	surname string,
	username string,
	birthDate time.Time,
	description string,
	email string,
	phoneNumber string,
	passwordHash string,
) *UserCreate {
	return &UserCreate{
		Name:         name,
		Surname:      surname,
		Username:     username,
		BirthDate:    birthDate,
		Description:  description,
		Email:        email,
		PhoneNumber:  phoneNumber,
		PasswordHash: passwordHash,
	}
}
