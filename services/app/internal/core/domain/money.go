package domain

import core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"

type Money struct {
	AmountMinorUnits int64
	Currency         string
}

func (m *Money) Validate() error {
	if m.AmountMinorUnits <= 0 {
		return core_errors.ErrInvalidBudget
	}
	if m.Currency != "RUB" {
		return core_errors.ErrUnsupportedCurrency
	}

	return nil
}
