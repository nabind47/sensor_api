package router

import (
	"net/http"

	"github.com/nabind47/sensor_api/internal/config"
	"github.com/nabind47/sensor_api/internal/handler"
	"github.com/nabind47/sensor_api/internal/middleware"
	"github.com/nabind47/sensor_api/internal/service"
	"github.com/nabind47/sensor_api/internal/storage"
)

func New(cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	repo := storage.NewMemoryStore()
	// repo := storage.NewFakeMemoryStore()
	srv := service.NewTemperatureService(repo)
	authMW := middleware.NewAuth(&cfg.Auth)

	handler := handler.NewSensorHandler(srv)

	mux.HandleFunc("GET /temperature", handler.GetSensors)
	mux.Handle("POST /temperature", authMW.AuthMiddleware(http.HandlerFunc(handler.CreateSensor)))

	return mux
}
