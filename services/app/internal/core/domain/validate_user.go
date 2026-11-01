package domain

import (
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/nyaruka/phonenumbers"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func validateUserName(name string) error {
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

func validateUserSurname(surname string) error {
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

func validateUserUsername(username string) error {
	userUsernameLength := len([]rune(strings.TrimSpace(username)))
	if userUsernameLength < 1 || userUsernameLength > 100 {
		return core_errors.ErrInvalidUserUsername
	}
	if strings.Contains(username, " ") {
		return core_errors.ErrInvalidUserUsername
	}
	return nil
}

func validateUserBirthDate(birthDate time.Time) error {
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

func validateUserDescription(description string) error {
	if len([]rune(description)) > 1000 {
		return core_errors.ErrInvalidUserDescription
	}
	return nil
}

func validateUserEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return core_errors.ErrInvalidUserEmail
	}
	return nil
}

func validateUserPhoneNumber(phoneNumber string) error {
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

func validateUserPassword(password string) error {
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

func validateUserDatesOrder(
	birthDate,
	createdAt,
	updatedAt time.Time,
	deletedAt *time.Time,
) error {
	if !birthDate.Before(createdAt) {
		return core_errors.ErrInvalidCreatedAt
	}

	return validateSystemDatesOrder(createdAt, updatedAt, deletedAt)
}
