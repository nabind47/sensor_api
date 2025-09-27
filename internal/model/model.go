package model

import (
	"errors"
	"time"
)

type SensorReading struct {
	SensorID    string    `json:"sensor_id"`
	Timestamp   time.Time `json:"timestamp"`
	Temperature float64   `json:"temperature"`
}

type SensorSummary struct {
	OverallAverage float64            `json:"overall_average"`
	SensorAverage  map[string]float64 `json:"sensor_average"`
}

type SensorRequestBody struct {
	SensorID    string   `json:"sensor_id"`
	Temperature *float64 `json:"temperature"`
	Timestamp   string   `json:"timestamp"`
}

func (n SensorRequestBody) Validate() (sensor SensorReading, err error) {
	if n.SensorID == "" {
		return sensor, errors.New("sensor_id is required")
	}

	if n.Timestamp == "" {
		return sensor, errors.New("timestamp is required")
	}

	if n.Temperature == nil {
		return sensor, errors.New("temperature is required")
	}

	temp := *n.Temperature
	if temp < -50 || temp > 100 {
		return sensor, errors.New("temperature out of valid range")
	}

	timestamp, parseErr := time.Parse(time.RFC3339, n.Timestamp)
	if parseErr != nil || timestamp.IsZero() {
		return sensor, errors.New("timestamp must be valid")
	}

	if timestamp.After(time.Now()) {
		return sensor, errors.New("timestamp cannot be in the future")
	}

	return SensorReading{
		SensorID:    n.SensorID,
		Temperature: temp,
		Timestamp:   timestamp,
	}, nil
}
