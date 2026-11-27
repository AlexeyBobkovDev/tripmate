package routes_service

import (
	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func routeDomainFromDTO(generatedRoute *route_generator_api_v1.GenerateRouteResponse) *domain.RouteCreate {
	return &domain.RouteCreate{
		Status: statusDomainFromDTO(generatedRoute.Status),
		Days:   daysDomainFromDTO(generatedRoute.Days),
	}
}

func statusDomainFromDTO(status route_generator_api_v1.RouteStatus) domain.RouteStatus {
	mapStatus := map[route_generator_api_v1.RouteStatus]domain.RouteStatus{
		route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED: domain.RouteStatusGenerated,
		route_generator_api_v1.RouteStatus_ROUTE_STATUS_EDITED:    domain.RouteStatusEdited,
	}

	mappedStatus, _ := mapStatus[status]

	return mappedStatus
}

func daysDomainFromDTO(days []*route_generator_api_v1.RouteDay) []domain.RouteDay {
	convertedDays := make([]domain.RouteDay, len(days))
	for ind, day := range days {
		convertedDays[ind] = domain.RouteDay{
			ID:        int(day.Id),
			Date:      day.Date.AsTime(),
			DayNumber: int(day.DayNumber),
			Stops:     stopsDomainFromDTO(day.Stops),
			Legs:      legsDomainFromDTO(day.Legs),
		}
	}
	return convertedDays
}

func legsDomainFromDTO(legs []*route_generator_api_v1.RouteLeg) []domain.RouteLeg {
	convertedDays := make([]domain.RouteLeg, len(legs))
	for ind, leg := range legs {
		convertedDays[ind] = domain.RouteLeg{
			ID:         int(leg.Id),
			FromStopID: int(leg.FromStopId),
			ToStopID:   int(leg.ToStopId),
			TravelMode: travelModeDTOToDomain(leg.TravelMode),
			Distance:   int(leg.Distance),
			Duration:   leg.Duration.AsDuration(),
		}
	}
	return convertedDays
}

func travelModeDTOToDomain(travelMode route_generator_api_v1.TravelMode) domain.TravelMode {
	mapTravelModes := map[route_generator_api_v1.TravelMode]domain.TravelMode{
		route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING: domain.TravelModeWalking,
		route_generator_api_v1.TravelMode_TRAVEL_MODE_CYCLING: domain.TravelModeCycling,
		route_generator_api_v1.TravelMode_TRAVEL_MODE_DRIVING: domain.TravelModeDriving,
		route_generator_api_v1.TravelMode_TRAVEL_MODE_TRANSIT: domain.TravelModeTransit,
	}
	mappedTravelMode, _ := mapTravelModes[travelMode]

	return mappedTravelMode
}

func stopsDomainFromDTO(stops []*route_generator_api_v1.RouteStop) []domain.RouteStop {
	convertedDays := make([]domain.RouteStop, len(stops))
	for ind, stop := range stops {
		convertedDays[ind] = domain.RouteStop{
			PlaceID:              int(stop.PlaceId),
			Position:             int(stop.Position),
			Name:                 stop.Name,
			Description:          stop.Description,
			Categories:           categoriesDomainFromDTO(stop.Categories),
			Interests:            interestsDomainFromDTO(stop.Interests),
			Address:              stop.Address,
			WebsiteURL:           stop.WebsiteUrl,
			Phone:                stop.Phone,
			VisitDurationMinutes: int(stop.VisitDurationMinutes),
			OpeningHours:         openingHoursDomainFromDTO(stop.OpeningHours),
			Coordinates: domain.Coordinates{
				Latitude:  stop.Coordinates.Latitude,
				Longitude: stop.Coordinates.Longitude,
			},
			ArrivalTime:   stop.ArrivalTime.AsTime(),
			DepartureTime: stop.DepartureTime.AsTime(),
		}
	}
	return convertedDays
}

func interestsDomainFromDTO(interests []route_generator_api_v1.Interest) []domain.Interest {
	mapInterests := map[route_generator_api_v1.Interest]domain.Interest{
		route_generator_api_v1.Interest_INTEREST_HISTORY:      domain.InterestHistory,
		route_generator_api_v1.Interest_INTEREST_NATURE:       domain.InterestNature,
		route_generator_api_v1.Interest_INTEREST_FOOD:         domain.InterestFood,
		route_generator_api_v1.Interest_INTEREST_ARCHITECTURE: domain.InterestArchitecture,
	}
	convertedInterests := make([]domain.Interest, len(interests))
	for ind, interest := range interests {
		mappedInterest, _ := mapInterests[interest]
		convertedInterests[ind] = mappedInterest
	}
	return convertedInterests
}

func categoriesDomainFromDTO(categories []route_generator_api_v1.Category) []domain.Category {
	mapCategories := map[route_generator_api_v1.Category]domain.Category{
		route_generator_api_v1.Category_CATEGORY_MUSEUM:        domain.CategoryMuseum,
		route_generator_api_v1.Category_CATEGORY_ART_GALLERY:   domain.CategoryArtGallery,
		route_generator_api_v1.Category_CATEGORY_THEATER:       domain.CategoryTheater,
		route_generator_api_v1.Category_CATEGORY_CINEMA:        domain.CategoryCinema,
		route_generator_api_v1.Category_CATEGORY_MONUMENT:      domain.CategoryMonument,
		route_generator_api_v1.Category_CATEGORY_HISTORICAL:    domain.CategoryHistorical,
		route_generator_api_v1.Category_CATEGORY_RELIGIOUS:     domain.CategoryReligious,
		route_generator_api_v1.Category_CATEGORY_PARK:          domain.CategoryPark,
		route_generator_api_v1.Category_CATEGORY_GARDEN:        domain.CategoryGarden,
		route_generator_api_v1.Category_CATEGORY_SQUARE:        domain.CategorySquare,
		route_generator_api_v1.Category_CATEGORY_STREET:        domain.CategoryStreet,
		route_generator_api_v1.Category_CATEGORY_ARCHITECTURE:  domain.CategoryArchitecture,
		route_generator_api_v1.Category_CATEGORY_OBSERVATION:   domain.CategoryObservation,
		route_generator_api_v1.Category_CATEGORY_ENTERTAINMENT: domain.CategoryEntertainment,
		route_generator_api_v1.Category_CATEGORY_RESTAURANT:    domain.CategoryRestaurant,
		route_generator_api_v1.Category_CATEGORY_CAFE:          domain.CategoryCafe,
		route_generator_api_v1.Category_CATEGORY_BAR:           domain.CategoryBar,
		route_generator_api_v1.Category_CATEGORY_SHOPPING:      domain.CategoryShopping,
		route_generator_api_v1.Category_CATEGORY_MARKET:        domain.CategoryMarket,
		route_generator_api_v1.Category_CATEGORY_NATURE:        domain.CategoryNature,
		route_generator_api_v1.Category_CATEGORY_ZOO:           domain.CategoryZoo,
		route_generator_api_v1.Category_CATEGORY_AQUARIUM:      domain.CategoryAquarium,
		route_generator_api_v1.Category_CATEGORY_SPORTS:        domain.CategorySports,
		route_generator_api_v1.Category_CATEGORY_EVENT_VENUE:   domain.CategoryEventVenue,
	}
	convertedCategories := make([]domain.Category, len(categories))
	for ind, category := range categories {
		mappedCategory, _ := mapCategories[category]
		convertedCategories[ind] = mappedCategory
	}
	return convertedCategories
}

func openingHoursDomainFromDTO(openingHours []*route_generator_common_v1.OpeningHours) []domain.OpeningHours {
	convertedOpeningHours := make([]domain.OpeningHours, len(openingHours))
	for ind := range openingHours {
		convertedOpeningHours[ind] = domain.OpeningHours{
			OpensAt:  int(openingHours[ind].OpensAt),
			ClosesAt: int(openingHours[ind].ClosesAt),
		}
	}
	return convertedOpeningHours
}
