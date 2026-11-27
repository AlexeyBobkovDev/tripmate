package core_http_request

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type innerDTO struct {
	City string `json:"city" validate:"required,min=2"`
}

type testDTO struct {
	Name  string   `json:"name" validate:"required,min=2,max=32"`
	Age   int      `json:"age" validate:"required,min=1,max=120"`
	Tags  []string `json:"tags" validate:"omitempty,max=3,dive,required"`
	Inner innerDTO `json:"inner" validate:"required"`
}

func (d *testDTO) Validate() error {
	if d.Name == "forbidden" {
		return core_errors.ErrInvalidTripName
	}

	return nil
}

func decodeBody(t *testing.T, body string) error {
	t.Helper()

	req := httptest.NewRequest("POST", "/test", strings.NewReader(body))

	return DecodeAndValidate(req, &testDTO{})
}

func TestDecodeAndValidate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		err := decodeBody(t, `{"name":"trip","age":30,"tags":["a"],"inner":{"city":"Moscow"}}`)
		require.NoError(t, err)
	})

	t.Run("empty body", func(t *testing.T) {
		err := decodeBody(t, "")
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("malformed json", func(t *testing.T) {
		err := decodeBody(t, `{"name":"trip",`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("not a json object - array", func(t *testing.T) {
		err := decodeBody(t, `["name","trip"]`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("not a json object - string", func(t *testing.T) {
		err := decodeBody(t, `"trip"`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("json null", func(t *testing.T) {
		err := decodeBody(t, `null`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("wrong type for number field", func(t *testing.T) {
		err := decodeBody(t, `{"name":"trip","age":"thirty","inner":{"city":"Moscow"}}`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("missing required field", func(t *testing.T) {
		err := decodeBody(t, `{"name":"trip","inner":{"city":"Moscow"}}`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("missing required nested field", func(t *testing.T) {
		err := decodeBody(t, `{"name":"trip","age":30,"inner":{}}`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("null nested object", func(t *testing.T) {
		err := decodeBody(t, `{"name":"trip","age":30,"inner":null}`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("custom Validate error", func(t *testing.T) {
		err := decodeBody(t, `{"name":"forbidden","age":30,"inner":{"city":"Moscow"}}`)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("unknown fields are ignored", func(t *testing.T) {
		err := decodeBody(t, `{"name":"trip","age":30,"inner":{"city":"Moscow"},"admin":true}`)
		require.NoError(t, err)
	})

	t.Run("duplicate json keys - last value wins", func(t *testing.T) {
		// encoding/json silently overwrites duplicate keys, so the last value
		// wins and the decoded destination reflects only the final value.
		req := httptest.NewRequest("POST", "/test", strings.NewReader(
			`{"name":"trip","age":30,"inner":{"city":"Moscow"},"age":30}`))
		dest := &testDTO{}
		require.NoError(t, DecodeAndValidate(req, dest))
		require.Equal(t, 30, dest.Age)

		req = httptest.NewRequest("POST", "/test", strings.NewReader(
			`{"name":"trip","age":30,"inner":{"city":"Moscow"},"age":5}`))
		dest = &testDTO{}
		require.NoError(t, DecodeAndValidate(req, dest))
		require.Equal(t, 5, dest.Age)
	})

	t.Run("trailing data after the json document is ignored", func(t *testing.T) {
		// BUG: json.Decoder stops after the first value, so any trailing
		// garbage after a complete JSON object is silently accepted.
		err := decodeBody(t, `{"name":"trip","age":30,"inner":{"city":"Moscow"}} {"age":-1}`)
		require.NoError(t, err)
	})
}

func TestValidate(t *testing.T) {
	t.Run("nil destination", func(t *testing.T) {
		err := Validate(nil)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("nil pointer to struct", func(t *testing.T) {
		var dto *testDTO
		err := Validate(dto)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("non-struct value", func(t *testing.T) {
		err := Validate("not a struct")
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("struct tag violation", func(t *testing.T) {
		err := Validate(&testDTO{Name: "a", Age: 30, Inner: innerDTO{City: "Moscow"}})
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})

	t.Run("slice element violation", func(t *testing.T) {
		err := Validate(&testDTO{Name: "trip", Age: 30, Tags: []string{""}, Inner: innerDTO{City: "Moscow"}})
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	})
}

func TestValidate_DiscardsSpecificError(t *testing.T) {
	// BUG: when the payload's own Validate() returns a specific domain error,
	// DecodeAndValidate replaces it with ErrInvalidArgument; the specific
	// cause is not wrapped, so callers cannot use errors.Is on it and API
	// clients cannot tell why the request failed.
	err := Validate(&testDTO{Name: "forbidden", Age: 30, Inner: innerDTO{City: "Moscow"}})

	require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
	require.False(t, errors.Is(err, core_errors.ErrInvalidTripName),
		"the specific domain error should be propagated to the caller but is discarded")
}
