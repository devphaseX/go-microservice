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

func (app *Config) errorJSON(w http.ResponseWriter, err error) error {
	errResponse := FromError(err)
	return writeJson(w, errResponse.Status, errResponse)
}

type SuccessApi struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta,omitempty"`
}

func NewSuccessApi(status int, data any, meta any, message string) SuccessApi {
	return SuccessApi{
		Status:  status,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func (app *Config) successJSON(w http.ResponseWriter, status int, data any, meta any, message string) error {
	resPayload := NewSuccessApi(status, data, meta, message)
	return writeJson(w, resPayload.Status, resPayload)
}
