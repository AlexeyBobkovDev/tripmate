package domain

type Interest string

const (
	InterestHistory      Interest = "history"
	InterestNature       Interest = "nature"
	InterestFood         Interest = "food"
	InterestArchitecture Interest = "architecture"
)

var AllowedInterests = map[Interest]struct{}{
	InterestHistory:      {},
	InterestNature:       {},
	InterestFood:         {},
	InterestArchitecture: {},
}
