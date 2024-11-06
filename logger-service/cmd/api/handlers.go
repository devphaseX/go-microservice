package main

import (
	"errors"
	"logger-service/data"
	"net/http"
)

type WriteLogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var payload WriteLogPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.errorJSON(w, NewError(errors.New("invalid bad request"), ErrBadRequest))
		return
	}

	event := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}

	if err := app.Models.LogEntity.Insert(&event); err != nil {
		app.errorJSON(w, NewError(errors.New("failed to write log"), ErrInternalFailure))
		return
	}

	app.successJSON(w, http.StatusCreated, nil, nil, "logged")
}
