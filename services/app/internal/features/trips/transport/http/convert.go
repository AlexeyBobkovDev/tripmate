package trips_transport_http

import (
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func tripDTOFromDomain(trip *domain.Trip) TripResponse {
	return TripResponse{
		ID:            trip.ID,
		OwnerID:       trip.OwnerID,
		Version:       trip.Version,
		Name:          trip.Name,
		PlanningArea:  trip.PlanningArea,
		StartLocation: trip.StartLocation,
		StartDate:     trip.StartDate,
		EndDate:       trip.EndDate,
		Budget:        trip.Budget,
		Interests:     trip.Interests,
		TimeZone:      trip.TimeZone,
		CreatedAt:     trip.CreatedAt,
		UpdatedAt:     trip.UpdatedAt,
		DeletedAt:     trip.DeletedAt,
	}
}

func tripDTOsFromDomains(trips *[]domain.Trip) []TripResponse {
	responses := make([]TripResponse, len(*trips))
	for ind, trip := range *trips {
		responses[ind] = tripDTOFromDomain(&trip)
	}

	return responses
}
