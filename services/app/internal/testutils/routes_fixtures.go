package testutils

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
)

func RouteGeneratorResponseTest(
	startDate time.Time,
) *route_generator_api_v1.GenerateRouteResponse {
	firstDayDate := time.Date(
		startDate.Year(),
		startDate.Month(),
		startDate.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	return &route_generator_api_v1.GenerateRouteResponse{
		Status: route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
		Days: []*route_generator_api_v1.RouteDay{
			routeGeneratorDayTest(1, firstDayDate, 100),
			routeGeneratorDayTest(2, firstDayDate.AddDate(0, 0, 1), 200),
		},
	}
}

func routeGeneratorDayTest(
	dayNumber int,
	date time.Time,
	stopIDBase int64,
) *route_generator_api_v1.RouteDay {
	firstStop := &route_generator_api_v1.RouteStop{
		Id:          stopIDBase,
		PlaceId:     stopIDBase + 1,
		Position:    1,
		Name:        fmt.Sprintf("Stop %d-%d", dayNumber, 1),
		Description: fmt.Sprintf("Description %d-%d", dayNumber, 1),
		Coordinates: &route_generator_common_v1.Coordinates{
			Latitude:  55.7539 + float64(dayNumber)/1000,
			Longitude: 37.6208 + float64(dayNumber)/1000,
		},
		ArrivalTime:   timestamppb.New(date.Add(9 * time.Hour)),
		DepartureTime: timestamppb.New(date.Add(11 * time.Hour)),
		VisitDuration: durationpb.New(90 * time.Minute),
	}

	secondStop := &route_generator_api_v1.RouteStop{
		Id:          stopIDBase + 1,
		PlaceId:     stopIDBase + 2,
		Position:    2,
		Name:        fmt.Sprintf("Stop %d-%d", dayNumber, 2),
		Description: fmt.Sprintf("Description %d-%d", dayNumber, 2),
		Coordinates: &route_generator_common_v1.Coordinates{
			Latitude:  55.7639 + float64(dayNumber)/1000,
			Longitude: 37.6308 + float64(dayNumber)/1000,
		},
		ArrivalTime:   timestamppb.New(date.Add(12 * time.Hour)),
		DepartureTime: timestamppb.New(date.Add(13*time.Hour + 30*time.Minute)),
		VisitDuration: durationpb.New(45 * time.Minute),
	}

	return &route_generator_api_v1.RouteDay{
		Id:        int64(dayNumber),
		Date:      timestamppb.New(date),
		DayNumber: int32(dayNumber),
		Stops:     []*route_generator_api_v1.RouteStop{firstStop, secondStop},
		Legs: []*route_generator_api_v1.RouteLeg{
			{
				Id:         int64(dayNumber),
				FromStopId: int32(firstStop.Id),
				ToStopId:   int32(secondStop.Id),
				TravelMode: route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING,
				Distance:   1500,
				Duration:   durationpb.New(15 * time.Minute),
			},
		},
	}
}
