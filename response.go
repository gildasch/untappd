package untappd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var (
	// errInvalidBool is returned when the Untappd API returns a
	// non 0 or 1 integer for a boolean value.
	errInvalidBool = errors.New("invalid boolean value")

	// errInvalidTimeUnit is returned when the Untappd API returns an
	// unrecognized time unit.
	errInvalidTimeUnit = errors.New("invalid time unit")
)

// responseTime implements json.Unmarshaler, so that duration responses
// in the Untappd APIv4 can be decoded directly into Go time.Duration structs.
type responseTime time.Duration

// UnmarshalJSON implements json.Unmarshaler.
func (r *responseTime) UnmarshalJSON(data []byte) error {
	var v struct {
		Time    float64 `json:"time"`
		Measure string  `json:"measure"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// Known measure strings mapped to Go parse-able equivalents
	timeUnits := map[string]string{
		"milliseconds": "ms",
		"seconds":      "s",
		"minutes":      "m",
	}

	// Parse a Go time.Duration from string
	d, err := time.ParseDuration(fmt.Sprintf("%f%s", v.Time, timeUnits[v.Measure]))
	if err != nil && strings.Contains(err.Error(), "time: missing unit in duration") {
		return errInvalidTimeUnit
	}

	*r = responseTime(d)
	return err
}

// responseURL implements json.Unmarshaler, so that URL string responses
// in the Untappd APIv4 can be decoded directly into Go *url.URL structs.
type responseURL url.URL

// UnmarshalJSON implements json.Unmarshaler.
func (r *responseURL) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	u, err := url.Parse(v)
	if err != nil {
		return err
	}

	*r = responseURL(*u)
	return nil
}

// responseBool implements json.Unmarshaler, so that integer 0 or 1 responses
// in the Untappd APIv4 can be decoded directly into Go boolean values.
type responseBool bool

// UnmarshalJSON implements json.Unmarshaler.
func (r *responseBool) UnmarshalJSON(data []byte) error {
	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch v {
	case 0:
		*r = false
	case 1:
		*r = true
	default:
		return errInvalidBool
	}

	return nil
}
