package domain

import (
	"math"
	"strings"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

func (l *Location) Validate() error {
	if len([]rune(strings.TrimSpace(l.Name))) == 0 {
		return core_errors.ErrInvalidLocationName
	}
	if math.IsNaN(l.Latitude) ||
		math.IsInf(l.Latitude, 0) ||
		l.Latitude < -90 ||
		l.Latitude > 90 {
		return core_errors.ErrInvalidLatitude
	}
	if math.IsNaN(l.Longitude) ||
		math.IsInf(l.Longitude, 0) ||
		l.Longitude < -180 ||
		l.Longitude > 180 {
		return core_errors.ErrInvalidLongitude
	}

	return nil
}
