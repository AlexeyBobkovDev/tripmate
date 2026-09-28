package domain

import (
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/nyaruka/phonenumbers"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func validateID(id int) error {
	if id <= 0 {
		return core_errors.ErrInvalidUserID
	}

	return nil
}

func validateVersion(version int) error {
	if version <= 0 {
		return core_errors.ErrInvalidUserVersion
	}

	return nil
}

func validateName(name string) error {
	userNameLength := len([]rune(strings.TrimSpace(name)))
	if userNameLength < 1 || userNameLength > 50 {
		return core_errors.ErrInvalidUserName
	}
	for _, letter := range name {
		if !unicode.IsLetter(letter) {
			return core_errors.ErrInvalidUserName
		}
	}
	return nil
}

func validateSurname(surname string) error {
	userSurnameLength := len([]rune(strings.TrimSpace(surname)))
	if userSurnameLength < 1 || userSurnameLength > 50 {
		return core_errors.ErrInvalidUserSurname
	}
	for _, letter := range surname {
		if !unicode.IsLetter(letter) {
			return core_errors.ErrInvalidUserSurname
		}
	}
	return nil
}

func validateUserName(username string) error {
	userUsernameLength := len([]rune(strings.TrimSpace(username)))
	if userUsernameLength < 1 || userUsernameLength > 100 {
		return core_errors.ErrInvalidUserUsername
	}
	return nil
}

func validateBirthDate(birthDate time.Time) error {
	if birthDate.After(time.Now()) {
		return core_errors.ErrInvalidUserBirthdate
	}

	age := calculateCurrentAge(birthDate)
	if age < 14 {
		return core_errors.ErrInvalidUserBirthdate
	}
	if age > 150 {
		return core_errors.ErrInvalidUserBirthdate
	}
	return nil
}

func validateDescription(description string) error {
	if len([]rune(description)) > 1000 {
		return core_errors.ErrInvalidUserDescription
	}
	return nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return core_errors.ErrInvalidUserEmail
	}
	return nil
}

func validatePhoneNumber(phoneNumber string) error {
	num, err := phonenumbers.Parse(phoneNumber, "")
	if err != nil {
		return core_errors.ErrInvalidUserPhoneNumber
	}

	if !phonenumbers.IsValidNumber(num) {
		return core_errors.ErrInvalidUserPhoneNumber
	}

	if phoneNumber != phonenumbers.Format(num, phonenumbers.E164) {
		return core_errors.ErrInvalidUserPhoneNumber
	}
	return nil
}

func validatePassword(password string) error {
	if len(strings.TrimSpace(password)) < 8 {
		return core_errors.ErrInvalidUserPassword
	}
	digits, letters, specSymbol, upper, lower := 0, 0, 0, 0, 0
	for _, symbol := range password {
		switch {
		case unicode.IsDigit(symbol):
			digits++
		case unicode.IsLetter(symbol):
			letters++
			switch {
			case unicode.IsUpper(symbol):
				upper++
			case unicode.IsLower(symbol):
				lower++
			}
		case unicode.IsPunct(symbol):
			specSymbol++

		}
	}

	if digits < 1 || letters < 1 || specSymbol < 1 || upper < 1 || lower < 1 {
		return core_errors.ErrInvalidUserPassword
	}
	return nil
}

func validateCreatedAt(createdAt time.Time) error {
	if createdAt.After(time.Now()) {
		return core_errors.ErrInvalidUserCreatedAt
	}
	return nil
}

func validateUpdatedAt(updatedAt time.Time) error {
	if updatedAt.After(time.Now()) {
		return core_errors.ErrInvalidUserUpdatedAt
	}
	return nil
}

func validateDeletedAt(deletedAt *time.Time) error {
	if deletedAt == nil {
		return nil
	}
	if deletedAt.After(time.Now()) {
		return core_errors.ErrInvalidUserDeletedAt
	}
	return nil
}

func validateDatesOrder(
	birthDate,
	createdAt,
	updatedAt time.Time,
	deletedAt *time.Time,
) error {
	if !birthDate.Before(createdAt) {
		return core_errors.ErrInvalidUserCreatedAt
	}

	if createdAt.After(updatedAt) {
		return core_errors.ErrInvalidUserUpdatedAt
	}

	if deletedAt != nil && updatedAt.After(*deletedAt) {
		return core_errors.ErrInvalidUserDeletedAt
	}

	return nil
}
