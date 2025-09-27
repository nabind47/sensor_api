package storage

import "github.com/nabind47/sensor_api/internal/model"

type StoreInterface interface {
	Create(reading model.SensorReading) (model.SensorReading, error)
	Get() map[string][]model.SensorReading
	GetSummary() map[string]any
}
