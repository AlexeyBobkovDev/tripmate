package domain

import (
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func validateID(id int) error {
	if id <= 0 {
		return core_errors.ErrInvalidID
	}

	return nil
}

func validateVersion(version int) error {
	if version <= 0 {
		return core_errors.ErrInvalidVersion
	}

	return nil
}

func validateCreatedAt(createdAt time.Time) error {
	if createdAt.IsZero() {
		return core_errors.ErrInvalidCreatedAt
	}
	if createdAt.After(time.Now()) {
		return core_errors.ErrInvalidCreatedAt
	}
	return nil
}

func validateUpdatedAt(updatedAt time.Time) error {
	if updatedAt.IsZero() {
		return core_errors.ErrInvalidUpdatedAt
	}
	if updatedAt.After(time.Now()) {
		return core_errors.ErrInvalidUpdatedAt
	}
	return nil
}

func validateDeletedAt(deletedAt *time.Time) error {
	if deletedAt == nil {
		return nil
	}
	if deletedAt.After(time.Now()) {
		return core_errors.ErrInvalidDeletedAt
	}
	return nil
}

func validateSystemDatesOrder(
	createdAt,
	updatedAt time.Time,
	deletedAt *time.Time,
) error {
	if createdAt.After(updatedAt) {
		return core_errors.ErrInvalidUpdatedAt
	}

	if deletedAt != nil && updatedAt.After(*deletedAt) {
		return core_errors.ErrInvalidDeletedAt
	}

	return nil
}
