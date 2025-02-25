package main

import (
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
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	log.Println(user.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid creds"), http.StatusBadRequest)
		log.Println("Authentication-Service Error2:", err)
		return
	}
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid creds"), http.StatusBadRequest)
		log.Println(requestPayload.Password)
		log.Println("Authentication-Service Error3:", err)
		return
	}

	payload := jsonResponse {
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}