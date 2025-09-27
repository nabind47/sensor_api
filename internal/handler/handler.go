package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nabind47/sensor_api/internal/logger"
	"github.com/nabind47/sensor_api/internal/model"
	"github.com/nabind47/sensor_api/internal/service"
	"github.com/nabind47/sensor_api/internal/util"
)

type SensorHandler struct {
	service *service.TemperatureService
}

func NewSensorHandler(service *service.TemperatureService) *SensorHandler {
	return &SensorHandler{
		service: service,
	}
}

func (h *SensorHandler) CreateSensor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var reading model.SensorRequestBody

	if err := json.NewDecoder(r.Body).Decode(&reading); err != nil {
		log.Error("failed to decode the request", "error", err)
		util.WriteError(w, http.StatusBadRequest, "payloads are empty")
		return
	}

	sensorReading, err := reading.Validate()
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.service.CreateReading(sensorReading)

	util.WriteSuccess(w, http.StatusCreated, nil)
}

func (h *SensorHandler) GetSensors(w http.ResponseWriter, r *http.Request) {
	readings := h.service.GetReadings()

	util.WriteSuccess(w, http.StatusOK, readings)
}
