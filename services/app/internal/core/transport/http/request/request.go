package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type validatable interface {
	Validate() error
}

var requestValidator = validator.New(validator.WithRequiredStructEnabled())

func Validate(obj any) error {
	errInvalidArgument := fmt.Errorf(
		"request validation: %v: %w",
		obj,
		core_errors.ErrInvalidArgument,
	)
	if err := requestValidator.Struct(obj); err != nil {
		return errInvalidArgument
	}

	if v, ok := obj.(validatable); ok {
		if err := v.Validate(); err != nil {
			return errInvalidArgument
		}
	}

	return nil
}

func DecodeAndValidate(req *http.Request, dest any) error {
	if err := json.NewDecoder(req.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			req,
			core_errors.ErrInvalidArgument,
		)
	}

	return Validate(dest)
}
