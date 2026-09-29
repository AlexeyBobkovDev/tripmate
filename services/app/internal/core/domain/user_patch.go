package domain

import "time"

type UserPatch struct {
	Version     int
	Name        Nullable[string]
	Surname     Nullable[string]
	Username    Nullable[string]
	BirthDate   Nullable[time.Time]
	Description Nullable[string]
	Email       Nullable[string]
	PhoneNumber Nullable[string]
}

func (u *UserPatch) Validate() error {
	if err := validateVersion(u.Version); err != nil {
		return err
	}
	if u.Name.Set {
		if err := validateName(*u.Name.Value); err != nil {
			return err
		}
	}
	if u.Surname.Set {
		if err := validateSurname(*u.Surname.Value); err != nil {
			return err
		}
	}
	if u.Username.Set {
		if err := validateUserName(*u.Username.Value); err != nil {
			return err
		}
	}
	if u.BirthDate.Set {
		if err := validateBirthDate(*u.BirthDate.Value); err != nil {
			return err
		}
	}
	if u.Description.Set {
		if err := validateDescription(*u.Description.Value); err != nil {
			return err
		}
	}
	if u.Email.Set {
		if err := validateEmail(*u.Email.Value); err != nil {
			return err
		}
	}
	if u.PhoneNumber.Set {
		if err := validatePhoneNumber(*u.PhoneNumber.Value); err != nil {
			return err
		}
	}

	return nil
}

func NewUserPatch(
	version int,
	name Nullable[string],
	surname Nullable[string],
	username Nullable[string],
	birthDate Nullable[time.Time],
	description Nullable[string],
	email Nullable[string],
	phoneNumber Nullable[string],
) *UserPatch {
	return &UserPatch{
		Version:     version,
		Name:        name,
		Surname:     surname,
		Username:    username,
		BirthDate:   birthDate,
		Description: description,
		Email:       email,
		PhoneNumber: phoneNumber,
	}
}
