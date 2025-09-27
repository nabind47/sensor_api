package util

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOK    = "ok"
	StatusError = "error"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	writeJson(w, status, ErrorResponse{Status: StatusError, Error: msg})
}

// Write a success response
func WriteSuccess(w http.ResponseWriter, status int, data any) {
	writeJson(w, status, SuccessResponse{Status: StatusOK, Data: data})
}
