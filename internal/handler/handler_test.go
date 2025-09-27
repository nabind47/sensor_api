package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/nabind47/sensor_api/internal/config"
	"github.com/nabind47/sensor_api/internal/handler"
	"github.com/nabind47/sensor_api/internal/middleware"
	"github.com/nabind47/sensor_api/internal/model"
	"github.com/nabind47/sensor_api/internal/service"
	"github.com/nabind47/sensor_api/internal/util"
)

type fakeRepo struct {
	readings map[string][]model.SensorReading
}

func (f *fakeRepo) Create(r model.SensorReading) (model.SensorReading, error) {
	if f.readings == nil {
		f.readings = make(map[string][]model.SensorReading)
	}
	f.readings[r.SensorID] = append(f.readings[r.SensorID], r)
	return r, nil
}

func (f *fakeRepo) Get() map[string][]model.SensorReading {
	return f.readings
}

func (f *fakeRepo) GetSummary() map[string]any {
	return map[string]any{
		"overall_average": "overallAverage",
		"sensor_average":  "sensorAverage",
	}
}

func TestCreateSensor(t *testing.T) {
	testCases := []struct {
		name           string
		body           string
		setAuth        bool
		invalidAuth    bool
		expectedStatus int
	}{
		{
			name:           "no auth header",
			body:           `{"sensor_id":"s1","temperature":25,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        false,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid auth header",
			body:           `{"sensor_id":"s1","temperature":25,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			invalidAuth:    true,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "empty body",
			body:           `{}`,
			setAuth:        true,
			invalidAuth:    false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing sensor_id",
			body:           `{"temperature":25,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty sensor_id",
			body:           `{"sensor_id":"","temperature":25,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing temperature",
			body:           `{"sensor_id":"s1","timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "null temperature",
			body:           `{"sensor_id":"s1","temperature":null,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing timestamp",
			body:           `{"sensor_id":"s1","temperature":25}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty timestamp",
			body:           `{"sensor_id":"s1","temperature":25,"timestamp":""}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid timestamp format",
			body:           `{"sensor_id":"s1","temperature":25,"timestamp":"invalid-date"}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "future timestamp",
			body:           `{"sensor_id":"s1","temperature":25,"timestamp":"2030-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "temperature too low",
			body:           `{"sensor_id":"s1","temperature":-60,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "temperature too high",
			body:           `{"sensor_id":"s1","temperature":150,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "temperature at lower bound",
			body:           `{"sensor_id":"s1","temperature":-50,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "temperature at upper bound",
			body:           `{"sensor_id":"s1","temperature":100,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "zero temperature",
			body:           `{"sensor_id":"s1","temperature":0,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "negative temperature in range",
			body:           `{"sensor_id":"s1","temperature":-25,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "float temperature",
			body:           `{"sensor_id":"s1","temperature":25.5,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json format",
			body:           `{"sensor_id":"s1","temperature":25,"timestamp":"2024-01-01T10:00:00Z"`,
			setAuth:        true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "special characters in sensor_id",
			body:           `{"sensor_id":"sensor@#$%^&*()","temperature":25,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "numeric sensor_id",
			body:           `{"sensor_id":"12345","temperature":25,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "valid reading with all fields",
			body:           `{"sensor_id":"s1","temperature":25,"timestamp":"2024-01-01T10:00:00Z"}`,
			setAuth:        true,
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeAuthCfg := &config.AuthConfig{
				ClientID:     "AUTH_CLIENT_ID",
				ClientSecret: "AUTH_CLIENT_SECRET",
				ClientExpiry: time.Duration(120) * time.Second,
			}

			repo := &fakeRepo{}
			svc := service.NewTemperatureService(repo)
			h := handler.NewSensorHandler(svc)

			req := httptest.NewRequest(http.MethodPost, "/temprature", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")

			if tc.setAuth {
				if tc.invalidAuth {
					req.Header.Set("x-authorization-key", "invalid-hash-token")
				} else {
					hash := util.GenerateHash(fakeAuthCfg.ClientID, fakeAuthCfg.ClientSecret)
					req.Header.Set("x-authorization-key", hash)
				}
			}

			w := httptest.NewRecorder()

			authMW := middleware.NewAuth(fakeAuthCfg)
			handlerWithAuth := authMW.AuthMiddleware(http.HandlerFunc(h.CreateSensor))
			handlerWithAuth.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
