package model_test

import (
	"testing"
	"time"

	"github.com/nabind47/sensor_api/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SensorRequestBody_Validate(t *testing.T) {
	validTemp := 25.0
	pastTime := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	futureTime := time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)

	type expectations struct {
		err    string
		sensor model.SensorReading
	}

	testCases := []struct {
		name         string
		req          model.SensorRequestBody
		expectations expectations
	}{
		{
			name: "missing sensor_id",
			req: model.SensorRequestBody{
				Temperature: &validTemp,
				Timestamp:   pastTime,
			},
			expectations: expectations{err: "sensor_id is required"},
		},
		{
			name: "missing timestamp",
			req: model.SensorRequestBody{
				SensorID:    "sensor-123",
				Temperature: &validTemp,
			},
			expectations: expectations{err: "timestamp is required"},
		},
		{
			name: "missing temperature",
			req: model.SensorRequestBody{
				SensorID:  "sensor-123",
				Timestamp: pastTime,
			},
			expectations: expectations{err: "temperature is required"},
		},
		{
			name: "temperature too low",
			req: model.SensorRequestBody{
				SensorID:    "sensor-123",
				Timestamp:   pastTime,
				Temperature: func() *float64 { v := -60.0; return &v }(),
			},
			expectations: expectations{err: "temperature out of valid range"},
		},
		{
			name: "temperature too high",
			req: model.SensorRequestBody{
				SensorID:    "sensor-123",
				Timestamp:   pastTime,
				Temperature: func() *float64 { v := 150.0; return &v }(),
			},
			expectations: expectations{err: "temperature out of valid range"},
		},
		{
			name: "invalid timestamp format",
			req: model.SensorRequestBody{
				SensorID:    "sensor-123",
				Timestamp:   "not-a-date",
				Temperature: &validTemp,
			},
			expectations: expectations{err: "timestamp must be valid"},
		},
		{
			name: "future timestamp",
			req: model.SensorRequestBody{
				SensorID:    "sensor-123",
				Timestamp:   futureTime,
				Temperature: &validTemp,
			},
			expectations: expectations{err: "timestamp cannot be in the future"},
		},
		{
			name: "valid request",
			req: model.SensorRequestBody{
				SensorID:    "sensor-123",
				Timestamp:   pastTime,
				Temperature: &validTemp,
			},
			expectations: expectations{
				sensor: model.SensorReading{
					SensorID:    "sensor-123",
					Temperature: validTemp,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sensor, err := tc.req.Validate()
			if tc.expectations.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectations.err)
			} else {
				assert.NoError(t, err)

				parsedTime, parseErr := time.Parse(time.RFC3339, tc.req.Timestamp)
				require.NoError(t, parseErr)
				tc.expectations.sensor.Timestamp = parsedTime

				assert.Equal(t, tc.expectations.sensor, sensor)
			}
		})
	}
}
