package domain

import (
	"fmt"
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

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

func (u *User) Validate() error {
	if err := validateID(u.ID); err != nil {
		return err
	}
	if err := validateVersion(u.Version); err != nil {
		return err
	}
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
	if err := validateCreatedAt(u.CreatedAt); err != nil {
		return err
	}
	if err := validateUpdatedAt(u.UpdatedAt); err != nil {
		return err
	}
	if err := validateDeletedAt(u.DeletedAt); err != nil {
		return err
	}
	if err := validateDatesOrder(u.BirthDate, u.CreatedAt, u.UpdatedAt, u.DeletedAt); err != nil {
		return err
	}

	return nil
}

func (u *User) ApplyPatch(patch *UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if tmp.Version != patch.Version {
		return core_errors.ErrDifferentVersions
	}

	tmp.Version++
	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}
	if patch.Surname.Set {
		tmp.Surname = *patch.Surname.Value
	}
	if patch.Username.Set {
		tmp.Username = *patch.Username.Value
	}
	if patch.BirthDate.Set {
		tmp.BirthDate = *patch.BirthDate.Value
	}
	if patch.Description.Set {
		tmp.Description = *patch.Description.Value
	}
	if patch.Email.Set {
		tmp.Email = *patch.Email.Value
	}
	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = *patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
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
