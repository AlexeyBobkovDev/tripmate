package routes_service

import "time"

func parseTimeString(timeStr string) time.Time {
	parsedTime, _ := time.Parse("15:04", timeStr)
	return parsedTime
}
