package core_errors

import "errors"

var ErrInvalidBirthDateFormat = NewError(errors.New("invalid birth date format"), CategoryBadRequest)
