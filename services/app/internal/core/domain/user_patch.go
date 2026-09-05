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
	// panic("unimplemented")
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
