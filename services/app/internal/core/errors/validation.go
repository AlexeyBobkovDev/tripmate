package core_errors

import "errors"

var (
	ErrInvalidUserID          = NewError(errors.New("invalid user id"), CategoryBadRequest)
	ErrInvalidUserVersion     = NewError(errors.New("invalid user version"), CategoryBadRequest)
	ErrInvalidUserName        = NewError(errors.New("invalid user name"), CategoryBadRequest)
	ErrInvalidUserSurname     = NewError(errors.New("invalid user surname"), CategoryBadRequest)
	ErrInvalidUserUsername    = NewError(errors.New("invalid user username"), CategoryBadRequest)
	ErrInvalidUserBirthdate   = NewError(errors.New("invalid user birthdate"), CategoryBadRequest)
	ErrInvalidUserDescription = NewError(errors.New("invalid user description"), CategoryBadRequest)
	ErrInvalidUserEmail       = NewError(errors.New("invalid user email"), CategoryBadRequest)
	ErrInvalidUserPhoneNumber = NewError(errors.New("invalid user phone number"), CategoryBadRequest)
	ErrInvalidUserPassword    = NewError(errors.New("invalid user password"), CategoryBadRequest)
	ErrInvalidUserCreatedAt   = NewError(errors.New("invalid user createdAt"), CategoryBadRequest)
	ErrInvalidUserUpdatedAt   = NewError(errors.New("invalid user updatedAt"), CategoryBadRequest)
	ErrInvalidUserDeletedAt   = NewError(errors.New("invalid user DeletedAt"), CategoryBadRequest)
	ErrDifferentVersions      = NewError(errors.New("invalid patch version"), CategoryBadRequest)
)
