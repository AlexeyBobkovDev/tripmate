package domain

import "time"

func calculateCurrentAge(birthdate time.Time) int {
	return calculateAge(birthdate, time.Now())
}

func calculateAge(birthDate, now time.Time) int {
	age := now.Year() - birthDate.Year()

	birthDateInThisYear := time.Date(
		now.Year(),
		birthDate.Month(),
		birthDate.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	if now.Before(birthDateInThisYear) {
		age--
	}

	return age
}
