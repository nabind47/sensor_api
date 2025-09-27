package storage

import (
	"sync"

	"github.com/nabind47/sensor_api/internal/model"
)

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string][]model.SensorReading
}

func NewMemoryStore() StoreInterface {
	return &MemoryStore{data: make(map[string][]model.SensorReading)}
}

func (m *MemoryStore) Create(reading model.SensorReading) (model.SensorReading, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[reading.SensorID] = append(m.data[reading.SensorID], reading)
	return reading, nil
}

func (m *MemoryStore) Get() map[string][]model.SensorReading {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data
}

func (m *MemoryStore) GetSummary() map[string]any {
	return nil
}
