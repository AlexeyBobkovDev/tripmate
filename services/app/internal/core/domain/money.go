package domain

import core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"

type Money struct {
	AmountMinorUnits int64  `json:"amount_minor_units"    validate:"required,min=1"  example:"1500000"`
	Currency         string `json:"currency"              validate:"required,alpha"  example:"RUB"`
}

func (m Money) Validate() error {
	if m.AmountMinorUnits <= 0 {
		return core_errors.ErrInvalidBudget
	}
	if m.Currency != "RUB" {
		return core_errors.ErrUnsupportedCurrency
	}

	return nil
}
