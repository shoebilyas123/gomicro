package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload);

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest);
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email);

	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest);
	}

	valid, err := user.PasswordMatches(requestPayload.Password);

	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest);
	}

	err = app.logRequest("Authentication",fmt.Sprintf("%s %s logged in.",user.FirstName, user.LastName))

	if err != nil {
		app.errorJSON(w, err);
		return;
	}

	payload := JSONResponse{
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.FirstName),
		Data: user, 
	}

	app.writeJSON(w, http.StatusAccepted, payload);
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string
		Data string
	}

	entry.Name = name;
	entry.Data = data;

	jsonData, _ := json.MarshalIndent(entry, "","\t");

	loggerSVCURL := "http://logsvc/log";


	request, err  := http.NewRequest("POST",loggerSVCURL, bytes.NewBuffer(jsonData));

	if err != nil {
		return err;
	}

	cl := &http.Client{}

	_, err = cl.Do(request);

	if err != nil {
		return err;
	}

	return nil;
}