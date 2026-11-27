package domain

type Category string

const (
	CategoryMuseum        Category = "museum"
	CategoryArtGallery    Category = "art_gallery"
	CategoryTheater       Category = "theater"
	CategoryCinema        Category = "cinema"
	CategoryMonument      Category = "monument"
	CategoryHistorical    Category = "historical_site"
	CategoryReligious     Category = "religious_site"
	CategoryPark          Category = "park"
	CategoryGarden        Category = "garden"
	CategorySquare        Category = "square"
	CategoryStreet        Category = "street"
	CategoryArchitecture  Category = "architecture"
	CategoryObservation   Category = "observation_deck"
	CategoryEntertainment Category = "entertainment"
	CategoryRestaurant    Category = "restaurant"
	CategoryCafe          Category = "cafe"
	CategoryBar           Category = "bar"
	CategoryShopping      Category = "shopping"
	CategoryMarket        Category = "market"
	CategoryNature        Category = "nature"
	CategoryZoo           Category = "zoo"
	CategoryAquarium      Category = "aquarium"
	CategorySports        Category = "sports"
	CategoryEventVenue    Category = "event_venue"
)

var AllowedCategories = map[Category]struct{}{
	CategoryMuseum:        {},
	CategoryArtGallery:    {},
	CategoryTheater:       {},
	CategoryCinema:        {},
	CategoryMonument:      {},
	CategoryHistorical:    {},
	CategoryReligious:     {},
	CategoryPark:          {},
	CategoryGarden:        {},
	CategorySquare:        {},
	CategoryStreet:        {},
	CategoryArchitecture:  {},
	CategoryObservation:   {},
	CategoryEntertainment: {},
	CategoryRestaurant:    {},
	CategoryCafe:          {},
	CategoryBar:           {},
	CategoryShopping:      {},
	CategoryMarket:        {},
	CategoryNature:        {},
	CategoryZoo:           {},
	CategoryAquarium:      {},
	CategorySports:        {},
	CategoryEventVenue:    {},
}
