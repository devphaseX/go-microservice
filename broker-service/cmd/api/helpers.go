package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxReadByteLimit := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxReadByteLimit))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
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

func (app *Config) errorJSON(w http.ResponseWriter, err any, status ...int) error {
	statusCode := http.StatusInternalServerError

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true

	if errValue, ok := err.(error); ok {
		payload.Message = errValue.Error()
	} else {
		payload.Data = err
	}

	return writeJson(w, statusCode, payload)
}
