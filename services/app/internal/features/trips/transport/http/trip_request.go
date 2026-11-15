package trips_transport_http

import (
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

type MealRequirementRequestDTO struct {
	Type     domain.MealType   `json:"meal_type"  validate:"required,oneof=breakfast lunch dinner"  example:"breakfast"`
	Window   domain.TimeWindow `json:"window"`
	Duration int               `json:"duration_minutes"  validate:"required,min=1,max=1440" example:"60"`
	Budget   domain.Money      `json:"budget"`
}

func (mr *MealRequirementRequestDTO) ToDomain() domain.MealRequirement {
	return domain.MealRequirement{
		Type:     mr.Type,
		Window:   mr.Window,
		Duration: time.Duration(mr.Duration) * time.Minute,
		Budget:   mr.Budget,
	}
}
