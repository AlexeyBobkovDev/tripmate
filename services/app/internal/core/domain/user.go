package domain

import "time"

type User struct {
	ID          int
	Version     int
	Name        string
	Surname     string
	Username    string
	BirthDate   time.Time
	Description string
	Email       string
	PhoneNumber string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func NewUser(
	id int,
	version int,
	name string,
	surname string,
	username string,
	birthDate time.Time,
	description string,
	email string,
	phoneNumber string,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *User {
	return &User{
		ID:          id,
		Version:     version,
		Name:        name,
		Surname:     surname,
		Username:    username,
		BirthDate:   birthDate,
		Description: description,
		Email:       email,
		PhoneNumber: phoneNumber,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		DeletedAt:   deletedAt,
	}
}

func NewUserUninitialized(
	name string,
	surname string,
	username string,
	birthDate time.Time,
	description string,
	email string,
	phoneNumber string,
) *User {
	return &User{
		ID:          UninitializedID,
		Version:     UninitializedVersion,
		Name:        name,
		Surname:     surname,
		Username:    username,
		BirthDate:   birthDate,
		Description: description,
		Email:       email,
		PhoneNumber: phoneNumber,
		DeletedAt:   nil,
	}
}
