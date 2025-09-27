package storage

import (
	"sync"

	"github.com/nabind47/sensor_api/internal/model"
)

type FakeMemoryStore struct {
	mu      sync.RWMutex
	records map[string][]model.SensorReading
}

func NewFakeMemoryStore() StoreInterface {
	return &FakeMemoryStore{records: make(map[string][]model.SensorReading)}
}

func (m *FakeMemoryStore) Create(reading model.SensorReading) (model.SensorReading, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.records[reading.SensorID] = append(m.records[reading.SensorID], reading)
	return reading, nil
}

func (m *FakeMemoryStore) Get() map[string][]model.SensorReading {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.records
}
