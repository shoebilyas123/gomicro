package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := JSONResponse{
		Error: false,
		Message: "Hit the broker server :) !!!",
	}

	log.Printf("BROKER: %T\n", payload);
	_ = app.writeJSON(w, http.StatusOK, payload);
}

type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload);

	
	if err != nil {
		app.errorJSON(w, err);
		return;
	}

	switch requestPayload.Action {
		case "auth":
			app.authenticate(w, requestPayload.Auth);
		default:
			app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
		// Create some JSON that we will send to the auth m-svc
		jsonData, _ := json.MarshalIndent(a, "", "\t");
		
		// Call the auth-svc
		resBody := bytes.NewBuffer(jsonData);
		request, err := http.NewRequest("POST", "http://authsvc/auth", resBody);

		if err != nil {
			app.errorJSON(w, err);
			return;
		}

		client := &http.Client{}
		response, err := client.Do(request);
		
		if err != nil {
			log.Println("CHECK1")
			app.errorJSON(w, err);
			return;
		}
		defer response.Body.Close();

		
		// Verify the response from svc and write the appropriate response to the calling client
		if response.StatusCode == http.StatusUnauthorized {
			log.Println("CHECK2")
			app.errorJSON(w, errors.New("invalid credentials"))
			return;
			} else if response.StatusCode != http.StatusAccepted {
				
			log.Println("CHECK3")
			app.errorJSON(w, errors.New("error calling auth service"))
			return;
		}
		
		var jsonFromService JSONResponse
		err = json.NewDecoder(response.Body).Decode(&jsonFromService);
		
		if err != nil {

			log.Println("CHECK4")
			app.errorJSON(w, err);
			return;
		}
		
		if jsonFromService.Error {

			log.Println("CHECK5")
			app.errorJSON(w, err, http.StatusUnauthorized);
			return;
		}

		var payload JSONResponse

		payload.Error = false;
		payload.Message = "Authenticated"
		payload.Data = jsonFromService.Data

		app.writeJSON(w, http.StatusAccepted, payload);
}
