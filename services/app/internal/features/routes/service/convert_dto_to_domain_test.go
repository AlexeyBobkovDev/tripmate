package routes_service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func TestConvertDtoToDomain_StatusDomainFromDTO(t *testing.T) {
	tests := []struct {
		name     string
		input    route_generator_api_v1.RouteStatus
		expected domain.RouteStatus
	}{
		{
			name:     "Generated status",
			input:    route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
			expected: domain.RouteStatusGenerated,
		},
		{
			name:     "Edited status",
			input:    route_generator_api_v1.RouteStatus_ROUTE_STATUS_EDITED,
			expected: domain.RouteStatusEdited,
		},
		{
			name:     "Unknown status",
			input:    route_generator_api_v1.RouteStatus_ROUTE_STATUS_UNSPECIFIED,
			expected: "", // zero value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := statusDomainFromDTO(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertDtoToDomain_TravelModeDTOToDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    route_generator_api_v1.TravelMode
		expected domain.TravelMode
	}{
		{
			name:     "Walking mode",
			input:    route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING,
			expected: domain.TravelModeWalking,
		},
		{
			name:     "Cycling mode",
			input:    route_generator_api_v1.TravelMode_TRAVEL_MODE_CYCLING,
			expected: domain.TravelModeCycling,
		},
		{
			name:     "Driving mode",
			input:    route_generator_api_v1.TravelMode_TRAVEL_MODE_DRIVING,
			expected: domain.TravelModeDriving,
		},
		{
			name:     "Transit mode",
			input:    route_generator_api_v1.TravelMode_TRAVEL_MODE_TRANSIT,
			expected: domain.TravelModeTransit,
		},
		{
			name:     "Unknown mode",
			input:    route_generator_api_v1.TravelMode_TRAVEL_MODE_UNSPECIFIED,
			expected: "", // zero value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := travelModeDTOToDomain(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertDtoToDomain_InterestsDomainFromDTO(t *testing.T) {
	tests := []struct {
		name     string
		input    []route_generator_api_v1.Interest
		expected []domain.Interest
	}{
		{
			name:     "Empty interests",
			input:    []route_generator_api_v1.Interest{},
			expected: []domain.Interest{},
		},
		{
			name:     "History interest",
			input:    []route_generator_api_v1.Interest{route_generator_api_v1.Interest_INTEREST_HISTORY},
			expected: []domain.Interest{domain.InterestHistory},
		},
		{
			name:     "Nature interest",
			input:    []route_generator_api_v1.Interest{route_generator_api_v1.Interest_INTEREST_NATURE},
			expected: []domain.Interest{domain.InterestNature},
		},
		{
			name:     "Food interest",
			input:    []route_generator_api_v1.Interest{route_generator_api_v1.Interest_INTEREST_FOOD},
			expected: []domain.Interest{domain.InterestFood},
		},
		{
			name:     "Architecture interest",
			input:    []route_generator_api_v1.Interest{route_generator_api_v1.Interest_INTEREST_ARCHITECTURE},
			expected: []domain.Interest{domain.InterestArchitecture},
		},
		{
			name:     "Unknown interest",
			input:    []route_generator_api_v1.Interest{route_generator_api_v1.Interest_INTEREST_UNSPECIFIED},
			expected: []domain.Interest{""}, // unknown interests are set to zero value
		},
		{
			name:     "Mixed interests",
			input:    []route_generator_api_v1.Interest{route_generator_api_v1.Interest_INTEREST_HISTORY, route_generator_api_v1.Interest_INTEREST_NATURE, route_generator_api_v1.Interest_INTEREST_UNSPECIFIED},
			expected: []domain.Interest{domain.InterestHistory, domain.InterestNature, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interestsDomainFromDTO(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertDtoToDomain_CategoriesDomainFromDTO(t *testing.T) {
	tests := []struct {
		name     string
		input    []route_generator_api_v1.Category
		expected []domain.Category
	}{
		{
			name:     "Empty categories",
			input:    []route_generator_api_v1.Category{},
			expected: []domain.Category{},
		},
		{
			name:     "Museum category",
			input:    []route_generator_api_v1.Category{route_generator_api_v1.Category_CATEGORY_MUSEUM},
			expected: []domain.Category{domain.CategoryMuseum},
		},
		{
			name:     "Park category",
			input:    []route_generator_api_v1.Category{route_generator_api_v1.Category_CATEGORY_PARK},
			expected: []domain.Category{domain.CategoryPark},
		},
		{
			name:     "Restaurant category",
			input:    []route_generator_api_v1.Category{route_generator_api_v1.Category_CATEGORY_RESTAURANT},
			expected: []domain.Category{domain.CategoryRestaurant},
		},
		{
			name:     "Unknown category",
			input:    []route_generator_api_v1.Category{route_generator_api_v1.Category_CATEGORY_UNSPECIFIED},
			expected: []domain.Category{""}, // unknown categories are set to zero value
		},
		{
			name:     "Mixed categories",
			input:    []route_generator_api_v1.Category{route_generator_api_v1.Category_CATEGORY_MUSEUM, route_generator_api_v1.Category_CATEGORY_PARK, route_generator_api_v1.Category_CATEGORY_UNSPECIFIED},
			expected: []domain.Category{domain.CategoryMuseum, domain.CategoryPark, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := categoriesDomainFromDTO(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertDtoToDomain_OpeningHoursDomainFromDTO(t *testing.T) {
	tests := []struct {
		name     string
		input    []*route_generator_common_v1.OpeningHours
		expected []domain.OpeningHours
	}{
		{
			name:     "Empty opening hours",
			input:    []*route_generator_common_v1.OpeningHours{},
			expected: []domain.OpeningHours{},
		},
		{
			name:     "Single opening hour",
			input:    []*route_generator_common_v1.OpeningHours{{OpensAt: 9, ClosesAt: 18}},
			expected: []domain.OpeningHours{{OpensAt: 9, ClosesAt: 18}},
		},
		{
			name:     "Multiple opening hours",
			input:    []*route_generator_common_v1.OpeningHours{{OpensAt: 9, ClosesAt: 18}, {OpensAt: 19, ClosesAt: 22}},
			expected: []domain.OpeningHours{{OpensAt: 9, ClosesAt: 18}, {OpensAt: 19, ClosesAt: 22}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := openingHoursDomainFromDTO(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertDtoToDomain_RouteDomainFromDTO(t *testing.T) {
	date := time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    *route_generator_api_v1.GenerateRouteResponse
		expected *domain.RouteCreate
	}{
		{
			name: "Empty route",
			input: &route_generator_api_v1.GenerateRouteResponse{
				Status: route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
				Days:   []*route_generator_api_v1.RouteDay{},
			},
			expected: &domain.RouteCreate{
				Status: domain.RouteStatusGenerated,
				Days:   []domain.RouteDay{},
			},
		},
		{
			name: "Route with stops",
			input: &route_generator_api_v1.GenerateRouteResponse{
				Status: route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
				Days: []*route_generator_api_v1.RouteDay{
					{
						Id:        1,
						Date:      timestamppb.New(date),
						DayNumber: 1,
						Stops: []*route_generator_api_v1.RouteStop{
							{
								PlaceId:              1,
								Position:             1,
								Name:                 "Test Stop",
								Description:          "Test Description",
								Address:              "Test Address",
								WebsiteUrl:           nil,
								Phone:                nil,
								VisitDurationMinutes: 60,
								Coordinates: &route_generator_common_v1.Coordinates{
									Latitude:  55.7539,
									Longitude: 37.6208,
								},
								ArrivalTime:   timestamppb.New(time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC)),
								DepartureTime: timestamppb.New(time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC)),
							},
						},
						Legs: []*route_generator_api_v1.RouteLeg{},
					},
				},
			},
			expected: &domain.RouteCreate{
				Status: domain.RouteStatusGenerated,
				Days: []domain.RouteDay{
					{
						ID:        1,
						Date:      date,
						DayNumber: 1,
						Stops: []domain.RouteStop{
							{
								ID:                   0,
								PlaceID:              1,
								Position:             1,
								Name:                 "Test Stop",
								Description:          "Test Description",
								Categories:           []domain.Category{},
								Interests:            []domain.Interest{},
								Address:              "Test Address",
								WebsiteURL:           nil,
								Phone:                nil,
								VisitDurationMinutes: 60,
								OpeningHours:         []domain.OpeningHours{},
								Coordinates: domain.Coordinates{
									Latitude:  55.7539,
									Longitude: 37.6208,
								},
								ArrivalTime:   time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC),
								DepartureTime: time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC),
							},
						},
						Legs: []domain.RouteLeg{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := routeDomainFromDTO(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
