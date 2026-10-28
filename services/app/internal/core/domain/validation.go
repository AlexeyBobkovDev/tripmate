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
		return core_errors.ErrInvalidID
	}

	return nil
}

func validateVersion(version int) error {
	if version <= 0 {
		return core_errors.ErrInvalidVersion
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
	if strings.Contains(username, " ") {
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
	if createdAt.IsZero() {
		return core_errors.ErrInvalidCreatedAt
	}
	if createdAt.After(time.Now()) {
		return core_errors.ErrInvalidCreatedAt
	}
	return nil
}

func validateUpdatedAt(updatedAt time.Time) error {
	if updatedAt.IsZero() {
		return core_errors.ErrInvalidUpdatedAt
	}
	if updatedAt.After(time.Now()) {
		return core_errors.ErrInvalidUpdatedAt
	}
	return nil
}

func validateDeletedAt(deletedAt *time.Time) error {
	if deletedAt == nil {
		return nil
	}
	if deletedAt.After(time.Now()) {
		return core_errors.ErrInvalidDeletedAt
	}
	return nil
}

func validateSystemDatesOrder(
	createdAt,
	updatedAt time.Time,
	deletedAt *time.Time,
) error {
	if createdAt.After(updatedAt) {
		return core_errors.ErrInvalidUpdatedAt
	}

	if deletedAt != nil && updatedAt.After(*deletedAt) {
		return core_errors.ErrInvalidDeletedAt
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

func validateStartDate(startDate time.Time) error {
	if startDate.IsZero() {
		return core_errors.ErrInvalidTripStartDate
	}
	return nil
}

func validateEndDate(endDate time.Time) error {
	if endDate.IsZero() {
		return core_errors.ErrInvalidTripEndDate
	}
	return nil
}

func validateTripDatesOrder(
	startDate,
	endDate time.Time,
	createdAt,
	updatedAt time.Time,
	deletedAt *time.Time,
) error {
	if !startDate.Before(endDate) {
		return core_errors.ErrInvalidTripEndDate
	}

	if startDate.Before(createdAt) || startDate.Equal(createdAt) {
		return core_errors.ErrInvalidTripStartDate
	}

	return validateSystemDatesOrder(createdAt, updatedAt, deletedAt)
}

func validateBudget(budget Money) error {
	return budget.Validate()
}

func validateTripName(name string) error {
	if len([]rune(strings.TrimSpace(name))) == 0 {
		return core_errors.ErrInvalidTripName
	}
	return nil
}

func validatePlanningArea(planningArea string) error {
	if len([]rune(strings.TrimSpace(planningArea))) == 0 {
		return core_errors.ErrInvalidPlanningArea
	}
	return nil
}

func validateStartLocation(startLocation Location) error {
	return startLocation.Validate()
}

func validateInterests(interests []Interest) error {
	for _, interest := range interests {
		if len([]rune(strings.TrimSpace(string(interest)))) == 0 {
			return core_errors.ErrInvalidInterest
		}
		if _, ok := AllowedInterests[interest]; !ok {
			return core_errors.ErrUnsupportedInterest
		}
	}
	return nil
}

func validateTimeZone(timeZone string) error {
	if len([]rune(strings.TrimSpace(timeZone))) == 0 {
		return core_errors.ErrInvalidTimeZone
	}
	if _, err := time.LoadLocation(timeZone); err != nil {
		return core_errors.ErrUnsupportedTimezone
	}
	if timeZone != "Europe/Moscow" {
		return core_errors.ErrUnsupportedTimezone
	}
	return nil
}
