package service

import (
	"github.com/nabind47/sensor_api/internal/model"
	"github.com/nabind47/sensor_api/internal/storage"
)

type TemperatureService struct {
	repo storage.StoreInterface
}

func NewTemperatureService(repo storage.StoreInterface) *TemperatureService {
	return &TemperatureService{repo: repo}
}

func (s *TemperatureService) CreateReading(reading model.SensorReading) (model.SensorReading, error) {
	return s.repo.Create(reading)
}

func (s *TemperatureService) GetReadings() model.SensorSummary {
	readings := s.repo.Get()

	return CalculateReadings(readings)
}

func CalculateReadings(readings map[string][]model.SensorReading) model.SensorSummary {
	summary := model.SensorSummary{
		SensorAverage:  make(map[string]float64),
		OverallAverage: 0,
	}

	var totalSum float64
	var totalCount int

	for sensorID, sensorReadings := range readings {
		var sum float64
		for _, r := range sensorReadings {
			sum += r.Temperature
		}

		if len(sensorReadings) > 0 {
			avg := sum / float64(len(sensorReadings))
			summary.SensorAverage[sensorID] = avg

			totalSum += sum
			totalCount += len(sensorReadings)
		}
	}

	if totalCount > 0 {
		summary.OverallAverage = totalSum / float64(totalCount)
	}

	return summary
}
