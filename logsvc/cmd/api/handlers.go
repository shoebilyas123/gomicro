package main

import (
	"logsvc/data"
	"net/http"
)


type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`

}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	_ = app.readJSON(w, r, &requestPayload);

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := event.InsertOne();
	
	if err != nil {
		app.errorJSON(w, err);
		return;
	}


	resp := JSONResponse{
		Error: false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp);

}