package domain

import (
	"math"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

type Coordinates struct {
	Latitude  float64
	Longitude float64
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
