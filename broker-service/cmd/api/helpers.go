package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func readJSON(w http.ResponseWriter, data any) error {
	maxReadByteLimit := 1048576

	reader := http.MaxBytesReader(w, nil, int64(maxReadByteLimit))

	return json.NewDecoder(reader).Decode(data)
}

func writeJson(w http.ResponseWriter, statusCode int, data any, headers ...http.Header) error {
	b, err := json.Marshal(data)

	if err != nil {
		return fmt.Errorf("invalid json payload: %w", err)
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(b)
	return err
}
