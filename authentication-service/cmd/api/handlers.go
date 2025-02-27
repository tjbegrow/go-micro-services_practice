package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate (w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println("Authentication-Service Error1:", err)
		return
	}
	user, err := app.Repo.GetByEmail(requestPayload.Email)
	log.Println(user.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid creds"), http.StatusBadRequest)
		log.Println("Authentication-Service Error2:", err)
		return
	}
	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid creds"), http.StatusBadRequest)
		log.Println(requestPayload.Password)
		log.Println("Authentication-Service Error3:", err)
		return
	}

	//log authentication with logger service
	//err = app.logRequest(w, "authentication", fmt.Sprintf("%s logged in", user.Email))
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	

	payload := jsonResponse {
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string`json:"name"`
		Data string`json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"
	
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	
	_, err = app.Client.Do(request)
	if err != nil {
		return err
	}
	return nil
}