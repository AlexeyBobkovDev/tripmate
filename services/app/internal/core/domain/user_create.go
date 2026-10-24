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
}

func (u *UserCreate) Validate() error {
	if err := validateName(u.Name); err != nil {
		return err
	}
	if err := validateSurname(u.Surname); err != nil {
		return err
	}
	if err := validateUserName(u.Username); err != nil {
		return err
	}
	if err := validateBirthDate(u.BirthDate); err != nil {
		return err
	}
	if err := validateDescription(u.Description); err != nil {
		return err
	}
	if err := validateEmail(u.Email); err != nil {
		return err
	}
	if err := validatePhoneNumber(u.PhoneNumber); err != nil {
		return err
	}

	return nil
}

func ValidatePassword(password string) error {
	return validatePassword(password)
}

func NewUserCreate(
	name string,
	surname string,
	username string,
	birthDate time.Time,
	description string,
	email string,
	phoneNumber string,
) *UserCreate {
	return &UserCreate{
		Name:        name,
		Surname:     surname,
		Username:    username,
		BirthDate:   birthDate,
		Description: description,
		Email:       email,
		PhoneNumber: phoneNumber,
	}
}
