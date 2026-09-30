package core_postgres_pool_errors

import (
	"errors"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

var (
	ErrNoRows             = core_errors.NewError(errors.New("no rows"), core_errors.CategoryDatabase)
	ErrViolatedForeignKey = core_errors.NewError(errors.New("violates foreign key"), core_errors.CategoryDatabase)
	ErrUnknown            = core_errors.NewError(errors.New("unknown"), core_errors.CategoryDatabase)
)
