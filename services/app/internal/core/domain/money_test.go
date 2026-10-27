package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestMoney_Validate(t *testing.T) {
	tests := []struct {
		name    string
		money   Money
		wantErr error
	}{
		{
			name: "valid minimum positive amount",
			money: Money{
				AmountMinorUnits: 1,
				Currency:         "RUB",
			},
		},
		{
			name: "valid typical amount",
			money: Money{
				AmountMinorUnits: 150000,
				Currency:         "RUB",
			},
		},
		{
			name: "valid maximum int64 amount",
			money: Money{
				AmountMinorUnits: math.MaxInt64,
				Currency:         "RUB",
			},
		},
		{
			name: "zero amount",
			money: Money{
				AmountMinorUnits: 0,
				Currency:         "RUB",
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "negative amount",
			money: Money{
				AmountMinorUnits: -1,
				Currency:         "RUB",
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "negative amount and unsupported currency",
			money: Money{
				AmountMinorUnits: -100,
				Currency:         "USD",
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "empty currency",
			money: Money{
				AmountMinorUnits: 100,
				Currency:         "",
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "unsupported USD currency",
			money: Money{
				AmountMinorUnits: 100,
				Currency:         "USD",
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "unsupported EUR currency",
			money: Money{
				AmountMinorUnits: 100,
				Currency:         "EUR",
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "lowercase currency",
			money: Money{
				AmountMinorUnits: 100,
				Currency:         "rub",
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "currency with leading whitespace",
			money: Money{
				AmountMinorUnits: 100,
				Currency:         " RUB",
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "currency with trailing whitespace",
			money: Money{
				AmountMinorUnits: 100,
				Currency:         "RUB ",
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "zero amount and empty currency",
			money: Money{
				AmountMinorUnits: 0,
				Currency:         "",
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.money.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
