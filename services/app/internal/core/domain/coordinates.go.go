package domain

import (
	"math"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"   validate:"required,latitude"        example:"55.753900"`
	Longitude float64 `json:"longitude"  validate:"required,longitude"       example:"37.620800"`
}

func (c Coordinates) Validate() error {
	if math.IsNaN(c.Latitude) ||
		math.IsInf(c.Latitude, 0) ||
		c.Latitude < -90 ||
		c.Latitude > 90 {
		return core_errors.ErrInvalidLatitude
	}
	if math.IsNaN(c.Longitude) ||
		math.IsInf(c.Longitude, 0) ||
		c.Longitude < -180 ||
		c.Longitude > 180 {
		return core_errors.ErrInvalidLongitude
	}

	return nil
}
