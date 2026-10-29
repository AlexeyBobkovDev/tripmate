package trips_transport_http

import (
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

type TripResponse struct {
	ID      int `json:"id" example:"1"`
	OwnerID int `json:"owner_id" example:"1"`
	Version int `json:"version" example:"1"`

	Name          string          `json:"name" example:"Walk around Moscow"`
	PlanningArea  string          `json:"planning_area" example:"Moscow"`
	StartLocation domain.Location `json:"start_location"`

	StartDate time.Time `json:"start_date" example:"2026-11-14T00:00:00Z"`
	EndDate   time.Time `json:"end_date" example:"2026-11-16T00:00:00Z"`

	Budget    domain.Money      `json:"budget"`
	Interests []domain.Interest `json:"interests"`
	TimeZone  string            `json:"time_zone" example:"Europe/Moscow"`

	CreatedAt time.Time  `json:"created_at" example:"2026-10-28T12:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2026-10-28T12:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" example:""`
}
