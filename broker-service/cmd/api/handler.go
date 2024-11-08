package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type RequestPayload struct {
	Action  string              `json:"action"`
	Auth    AuthPayload         `json:"auth"`
	Log     LogPayload          `json:"log"`
	NewUser RegisterUserPayload `json:"new_user"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Config) root(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Default root path",
	}

	writeJson(w, http.StatusAccepted, payload)
}

func (c *Config) handleSubmission(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload

	err := readJSON(w, r, &payload)

	if err != nil {
		c.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch payload.Action {
	case "auth.sign_in":
		c.authenicate(w, payload.Auth)

	case "auth.sign_up":
		c.registerUser(w, payload.NewUser)
	case "log":
		c.log(w, payload.Log)
	}
}

func (c *Config) authenicate(w http.ResponseWriter, p AuthPayload) {
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/auth/sign-in", bytes.NewBuffer(jsonData))

	if err != nil {
		c.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		c.errorJSON(w, err)
		return
	}

	defer response.Body.Close()
	var respPayload map[string]any
	json.NewDecoder(response.Body).Decode(&respPayload)

	if response.StatusCode != http.StatusOK {
		c.errorJSON(w, respPayload, response.StatusCode)
		return
	}

	writeJson(w, http.StatusOK, respPayload)
}

func (c *Config) registerUser(w http.ResponseWriter, p RegisterUserPayload) {
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/auth/sign-up", bytes.NewBuffer(jsonData))

	if err != nil {
		c.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		c.errorJSON(w, err)
		return
	}

	defer response.Body.Close()
	var respPayload map[string]any
	json.NewDecoder(response.Body).Decode(&respPayload)

	if response.StatusCode != http.StatusOK {
		c.errorJSON(w, respPayload, response.StatusCode)
		return
	}

	writeJson(w, http.StatusOK, respPayload)
}

func (c *Config) log(w http.ResponseWriter, p LogPayload) {
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	request, err := http.NewRequest(http.MethodPost, "http://logger-service", bytes.NewBuffer(jsonData))

	if err != nil {
		c.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		c.errorJSON(w, err)
		return
	}

	defer response.Body.Close()
	var respPayload map[string]any
	json.NewDecoder(response.Body).Decode(&respPayload)

	if response.StatusCode != http.StatusOK {
		c.errorJSON(w, respPayload)
		return
	}

	writeJson(w, http.StatusOK, respPayload)
}
