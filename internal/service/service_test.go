package service_test

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/nabind47/sensor_api/internal/model"
	"github.com/nabind47/sensor_api/internal/service"
)

func TestCalculateReadings(t *testing.T) {
	type expectations struct {
		overall float64
		avgs    map[string]float64
	}

	now := time.Now()

	testCases := []struct {
		name     string
		readings map[string][]model.SensorReading
		expect   expectations
	}{
		{
			name: "two sensors with two readings each",
			readings: map[string][]model.SensorReading{
				"1": {
					{SensorID: "1", Temperature: 10.0, Timestamp: now},
					{SensorID: "1", Temperature: 20.0, Timestamp: now},
				},
				"2": {
					{SensorID: "2", Temperature: 30.0, Timestamp: now},
					{SensorID: "2", Temperature: 40.0, Timestamp: now},
				},
			},
			expect: expectations{
				overall: 25.0,
				avgs: map[string]float64{
					"1": 15.0,
					"2": 35.0,
				},
			},
		},
		{
			name:     "empty readings",
			readings: map[string][]model.SensorReading{},
			expect: expectations{
				overall: 0.0,
				avgs:    map[string]float64{},
			},
		},
		{
			name: "single sensor single reading",
			readings: map[string][]model.SensorReading{
				"1": {
					{SensorID: "1", Temperature: 50.0, Timestamp: now},
				},
			},
			expect: expectations{
				overall: 50.0,
				avgs: map[string]float64{
					"1": 50.0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			summary := service.CalculateReadings(tc.readings)

			assert.Equal(t, tc.expect.overall, summary.OverallAverage)
			assert.Equal(t, tc.expect.avgs, summary.SensorAverage)
		})
	}
}
