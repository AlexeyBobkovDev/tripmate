package domain

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want error
	}{
		{name: "minimum invalid", id: 0, want: core_errors.ErrInvalidID},
		{name: "negative", id: -1, want: core_errors.ErrInvalidID},
		{name: "minimum integer", id: math.MinInt, want: core_errors.ErrInvalidID},
		{name: "one", id: 1},
		{name: "positive", id: 42},
		{name: "maximum integer", id: math.MaxInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateID(tt.id)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version int
		want    error
	}{
		{name: "zero", version: 0, want: core_errors.ErrInvalidVersion},
		{name: "negative", version: -1, want: core_errors.ErrInvalidVersion},
		{name: "minimum integer", version: math.MinInt, want: core_errors.ErrInvalidVersion},
		{name: "one", version: 1},
		{name: "positive", version: 42},
		{name: "maximum integer", version: math.MaxInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateVersion(tt.version)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateCreatedAt(t *testing.T) {
	now := time.Now()
	require.NoError(t, validateCreatedAt(now.Add(-time.Second)))
	require.NoError(t, validateCreatedAt(now))
	require.ErrorIs(t, validateCreatedAt(now.Add(time.Second)), core_errors.ErrInvalidCreatedAt)
	require.ErrorIs(t, validateCreatedAt(time.Time{}), core_errors.ErrInvalidCreatedAt)
}

func TestValidateUpdatedAt(t *testing.T) {
	now := time.Now()
	require.NoError(t, validateUpdatedAt(now.Add(-time.Second)))
	require.NoError(t, validateUpdatedAt(now))
	require.ErrorIs(t, validateUpdatedAt(now.Add(time.Second)), core_errors.ErrInvalidUpdatedAt)
	require.ErrorIs(t, validateUpdatedAt(time.Time{}), core_errors.ErrInvalidUpdatedAt)
}

func TestValidateDeletedAt(t *testing.T) {
	now := time.Now()
	require.NoError(t, validateDeletedAt(nil))
	require.NoError(t, validateDeletedAt(timePtr(now.Add(-time.Second))))
	require.NoError(t, validateDeletedAt(timePtr(now)))
	require.ErrorIs(t, validateDeletedAt(timePtr(now.Add(time.Second))), core_errors.ErrInvalidDeletedAt)
	require.NoError(t, validateDeletedAt(timePtr(time.Time{})))
}

func TestValidateSystemDatesOrder(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		created time.Time
		updated time.Time
		deleted *time.Time
		want    error
	}{
		{name: "strict chronological order", created: now, updated: now.Add(time.Hour)},
		{name: "created equals updated", created: now, updated: now},
		{name: "deleted later than updated", created: now, updated: now.Add(time.Hour), deleted: timePtr(now.Add(2 * time.Hour))},
		{name: "deleted equals updated", created: now, updated: now.Add(time.Hour), deleted: timePtr(now.Add(time.Hour))},
		{name: "created after updated", created: now, updated: now.Add(-time.Second), want: core_errors.ErrInvalidUpdatedAt},
		{name: "updated after deleted", created: now, updated: now.Add(time.Hour), deleted: timePtr(now.Add(time.Minute)), want: core_errors.ErrInvalidDeletedAt},
		{name: "no deleted timestamp", created: now, updated: now, deleted: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSystemDatesOrder(tt.created, tt.updated, tt.deleted)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
