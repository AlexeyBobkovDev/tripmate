package generate_route_service

import (
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

type FindCandidatesCriteria struct {
	CityName  string
	Interests []domain.Interest
	Date      time.Time
}
