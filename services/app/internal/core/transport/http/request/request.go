package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

type validatable interface {
	Validate() error
}

var requestValidator = validator.New(validator.WithRequiredStructEnabled())

func DecodeAndValidate(req *http.Request, dest any) error {
	if err := json.NewDecoder(req.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			req,
			core_errors.ErrInvalidArgument,
		)
	}

	var err error

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			dest,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
