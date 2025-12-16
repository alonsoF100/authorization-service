package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrFailedToDecode   = errors.New("failed to decode JSON")
	ErrFailedToValidate = errors.New("failed to validate request")
	ErrServer           = errors.New("damn, the server gaz up for nothing")
)

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Debug("Failed to encode", "data", data, "error", err)
	}
}
