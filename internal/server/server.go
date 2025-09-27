package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nabind47/sensor_api/internal/config"
	"github.com/nabind47/sensor_api/internal/logger"
	"github.com/nabind47/sensor_api/internal/router"
)

func NewServer(cfg *config.Config, log *slog.Logger) *http.Server {
	r := router.New(cfg)
	wrappedRouter := logger.AddLoggerMiddleware(log, logger.LogRequestMiddleware(r))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      wrappedRouter,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	return server
}
