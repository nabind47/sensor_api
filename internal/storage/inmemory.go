package storage

import (
	"sync"

	"github.com/nabind47/sensor_api/internal/model"
)

type Sensor struct {
	Sum     float64               `json:"sum"`
	Count   int                   `json:"count"`
	Records []model.SensorReading `json:"records"`
}

func (s *Sensor) Average() float64 {
	if s.Count == 0 {
		return 0
	}
	return s.Sum / float64(s.Count)
}

type InMemoryStore struct {
	mu      sync.RWMutex
	sensors map[string]*Sensor
}

func NewInMemoryStore() StoreInterface {
	return &InMemoryStore{
		sensors: make(map[string]*Sensor),
	}
}

func (m *InMemoryStore) Create(reading model.SensorReading) (model.SensorReading, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sensor, exists := m.sensors[reading.SensorID]
	if !exists {
		sensor = &Sensor{
			Sum:     0,
			Count:   0,
			Records: make([]model.SensorReading, 0),
		}
		m.sensors[reading.SensorID] = sensor
	}

	sensor.Records = append(sensor.Records, reading)

	sensor.Sum += reading.Temperature
	sensor.Count++

	return reading, nil
}
func (m *InMemoryStore) Get() map[string][]model.SensorReading {
	return nil
}

func (m *InMemoryStore) GetSummary() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sensorAverage := make(map[string]float64)
	var totalSum float64
	var totalCount int

	for sensorID, sensor := range m.sensors {
		sensorAverage[sensorID] = sensor.Average()

		totalSum += sensor.Sum
		totalCount += sensor.Count
	}

	var overallAverage float64
	if totalCount > 0 {
		overallAverage = totalSum / float64(totalCount)
	}

	return map[string]any{
		"overall_average": overallAverage,
		"sensor_average":  sensorAverage,
	}
}
