package core_errors

import "errors"

var (
	ErrInvalidArgument = NewError(errors.New("invalid argument"), CategoryBadRequest)
	ErrNotFound        = NewError(errors.New("not found"), CategoryNotFound)
	ErrConflict        = NewError(errors.New("conflict"), CategoryConflict)
	ErrBadRequest      = NewError(errors.New("bad request"), CategoryBadRequest)
)
