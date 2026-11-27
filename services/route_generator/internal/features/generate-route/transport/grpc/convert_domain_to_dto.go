package transport_grpc_generate_route

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

func convertRouteDomainToDTO(route *domain.Route) *route_generator_api_v1.GenerateRouteResponse {
	routeStatusMap := map[domain.RouteStatus]route_generator_api_v1.RouteStatus{
		domain.RouteStatusGenerated: route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
		domain.RouteStatusEdited:    route_generator_api_v1.RouteStatus_ROUTE_STATUS_EDITED,
	}
	routeStatus, _ := routeStatusMap[route.Status]

	return &route_generator_api_v1.GenerateRouteResponse{
		Status: routeStatus,
		Days:   convertDaysFromDomainToDTO(route.Days),
	}
}

func convertDaysFromDomainToDTO(days []domain.RouteDay) []*route_generator_api_v1.RouteDay {
	convertedDays := make([]*route_generator_api_v1.RouteDay, len(days))
	for ind, day := range days {
		convertedDays[ind] = convertDayFromDomainToDTO(day)
	}
	return convertedDays
}

func convertDayFromDomainToDTO(day domain.RouteDay) *route_generator_api_v1.RouteDay {
	return &route_generator_api_v1.RouteDay{
		Id:        int64(day.ID),
		Date:      timestamppb.New(day.Date),
		DayNumber: int32(day.DayNumber),
		Stops:     convertStopsFromDomainToDTO(day.Stops),
		Legs:      convertLegsFromDomainToDTO(day.Legs),
	}
}

func convertStopsFromDomainToDTO(stops []domain.RouteStop) []*route_generator_api_v1.RouteStop {
	convertedStop := make([]*route_generator_api_v1.RouteStop, len(stops))
	for ind, stop := range stops {
		convertedStop[ind] = convertStopFromDomainToDTO(stop)
	}
	return convertedStop
}

func convertStopFromDomainToDTO(stop domain.RouteStop) *route_generator_api_v1.RouteStop {
	var (
		websiteURL *string = nil
		phone      *string = nil
	)
	if stop.WebsiteURL != nil {
		websiteURL = stop.WebsiteURL
	}
	if stop.Phone != nil {
		phone = stop.Phone
	}
	return &route_generator_api_v1.RouteStop{
		PlaceId:              int64(stop.PlaceID),
		Position:             int32(stop.Position),
		Name:                 stop.Name,
		Description:          stop.Description,
		Categories:           categoriesDomainToDTO(stop.Categories),
		Interests:            interestsDomainToDTO(stop.Interests),
		Address:              stop.Address,
		WebsiteUrl:           websiteURL,
		Phone:                phone,
		VisitDurationMinutes: int32(stop.VisitDurationMinutes),
		OpeningHours:         openingHoursDomainToDTO(stop.OpeningHours),
		Coordinates:          convertStartCoordinatesDomainToDTO(stop.Coordinates),
		ArrivalTime:          timestamppb.New(stop.ArrivalTime),
		DepartureTime:        timestamppb.New(stop.DepartureTime),
	}
}

func interestsDomainToDTO(interests []domain.Interest) []route_generator_api_v1.Interest {
	mapInterests := map[domain.Interest]route_generator_api_v1.Interest{
		domain.InterestHistory:      route_generator_api_v1.Interest_INTEREST_HISTORY,
		domain.InterestNature:       route_generator_api_v1.Interest_INTEREST_NATURE,
		domain.InterestFood:         route_generator_api_v1.Interest_INTEREST_FOOD,
		domain.InterestArchitecture: route_generator_api_v1.Interest_INTEREST_ARCHITECTURE,
	}
	convertedInterests := make([]route_generator_api_v1.Interest, len(interests))
	for ind, interest := range interests {
		mappedInterest, _ := mapInterests[interest]
		convertedInterests[ind] = mappedInterest
	}
	return convertedInterests
}

func categoriesDomainToDTO(categories []domain.Category) []route_generator_api_v1.Category {
	mapCategories := map[domain.Category]route_generator_api_v1.Category{
		domain.CategoryMuseum:        route_generator_api_v1.Category_CATEGORY_MUSEUM,
		domain.CategoryArtGallery:    route_generator_api_v1.Category_CATEGORY_ART_GALLERY,
		domain.CategoryTheater:       route_generator_api_v1.Category_CATEGORY_THEATER,
		domain.CategoryCinema:        route_generator_api_v1.Category_CATEGORY_CINEMA,
		domain.CategoryMonument:      route_generator_api_v1.Category_CATEGORY_MONUMENT,
		domain.CategoryHistorical:    route_generator_api_v1.Category_CATEGORY_HISTORICAL,
		domain.CategoryReligious:     route_generator_api_v1.Category_CATEGORY_RELIGIOUS,
		domain.CategoryPark:          route_generator_api_v1.Category_CATEGORY_PARK,
		domain.CategoryGarden:        route_generator_api_v1.Category_CATEGORY_GARDEN,
		domain.CategorySquare:        route_generator_api_v1.Category_CATEGORY_SQUARE,
		domain.CategoryStreet:        route_generator_api_v1.Category_CATEGORY_STREET,
		domain.CategoryArchitecture:  route_generator_api_v1.Category_CATEGORY_ARCHITECTURE,
		domain.CategoryObservation:   route_generator_api_v1.Category_CATEGORY_OBSERVATION,
		domain.CategoryEntertainment: route_generator_api_v1.Category_CATEGORY_ENTERTAINMENT,
		domain.CategoryRestaurant:    route_generator_api_v1.Category_CATEGORY_RESTAURANT,
		domain.CategoryCafe:          route_generator_api_v1.Category_CATEGORY_CAFE,
		domain.CategoryBar:           route_generator_api_v1.Category_CATEGORY_BAR,
		domain.CategoryShopping:      route_generator_api_v1.Category_CATEGORY_SHOPPING,
		domain.CategoryMarket:        route_generator_api_v1.Category_CATEGORY_MARKET,
		domain.CategoryNature:        route_generator_api_v1.Category_CATEGORY_NATURE,
		domain.CategoryZoo:           route_generator_api_v1.Category_CATEGORY_ZOO,
		domain.CategoryAquarium:      route_generator_api_v1.Category_CATEGORY_AQUARIUM,
		domain.CategorySports:        route_generator_api_v1.Category_CATEGORY_SPORTS,
		domain.CategoryEventVenue:    route_generator_api_v1.Category_CATEGORY_EVENT_VENUE,
	}
	convertedCategories := make([]route_generator_api_v1.Category, len(categories))
	for ind, category := range categories {
		mappedCategory, _ := mapCategories[category]
		convertedCategories[ind] = mappedCategory
	}
	return convertedCategories
}

func openingHoursDomainToDTO(openingHours []domain.OpeningHours) []*route_generator_common_v1.OpeningHours {
	convertedOpeningHours := make([]*route_generator_common_v1.OpeningHours, len(openingHours))
	for ind := range openingHours {
		convertedOpeningHours[ind] = &route_generator_common_v1.OpeningHours{
			OpensAt:  int32(openingHours[ind].OpensAt),
			ClosesAt: int32(openingHours[ind].ClosesAt),
		}
	}
	return convertedOpeningHours
}

func convertLegsFromDomainToDTO(legs []domain.RouteLeg) []*route_generator_api_v1.RouteLeg {
	convertedLegs := make([]*route_generator_api_v1.RouteLeg, len(legs))
	for ind, leg := range legs {
		convertedLegs[ind] = convertLegFromDomainToDTO(leg)
	}
	return convertedLegs
}

func convertLegFromDomainToDTO(leg domain.RouteLeg) *route_generator_api_v1.RouteLeg {
	travelModeMap := map[domain.TravelMode]route_generator_api_v1.TravelMode{
		domain.TravelModeCycling: route_generator_api_v1.TravelMode_TRAVEL_MODE_CYCLING,
		domain.TravelModeWalking: route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING,
		domain.TravelModeDriving: route_generator_api_v1.TravelMode_TRAVEL_MODE_DRIVING,
		domain.TravelModeTransit: route_generator_api_v1.TravelMode_TRAVEL_MODE_TRANSIT,
	}
	travelMode, _ := travelModeMap[leg.TravelMode]

	return &route_generator_api_v1.RouteLeg{
		Id:         int64(leg.ID),
		FromStopId: int32(leg.FromStopID),
		ToStopId:   int32(leg.ToStopID),
		TravelMode: travelMode,
		Distance:   int32(leg.Distance),
		Duration:   durationpb.New(leg.Duration),
	}
}

func convertStartCoordinatesDomainToDTO(startCoordinates domain.Coordinates) *route_generator_common_v1.Coordinates {
	return &route_generator_common_v1.Coordinates{
		Latitude:  startCoordinates.Latitude,
		Longitude: startCoordinates.Longitude,
	}
}
