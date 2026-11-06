package routes_transport_http

import (
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

type RouteResponse struct {
	ID          int `json:"id" example:"1"`
	TripID      int `json:"owner_id" example:"1"`
	TripVersion int `json:"trip_version" example:"1"`
	Version     int `json:"version" example:"1"`

	Status domain.RouteStatus `json:"status" example:"generated"`
	Days   []domain.RouteDay  `json:"days"`

	CreatedAt time.Time `json:"created_at" example:"2026-10-28T12:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-10-28T12:00:00Z"`
}
