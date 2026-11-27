package domain

type Place struct {
	ID                   int
	CityName             string
	Name                 string
	Description          string
	ShortDescription     string
	PlaceType            PlaceType
	Categories           []string
	Interests            []string
	Latitude             float64
	Longitude            float64
	Address              string
	WebsiteURL           *string
	Phone                *string
	VisitDurationMinutes int
	Status               string
	OpeningHours         []Schedule
	Exceptions           []Exception
}
