package domain

import (
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func validateID(id int) error {
	if id <= 0 {
		return core_errors.ErrInvalidID
	}

	return nil
}
