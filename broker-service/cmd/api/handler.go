package main

import (
	"net/http"
)

func (c *Config) root(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Default root path",
	}

	writeJson(w, http.StatusAccepted, payload)
}
