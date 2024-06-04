package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// The limit of the JSON payload that the client can pass.
	maxBytes := 1048576;

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes));

	decoder := json.NewDecoder(r.Body); // Create a new decoder with r.Body as it's source
	err := decoder.Decode(data) // Decode the r.Body into the data.

	if err != nil {
		return err;
	}

	err = decoder.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must not contain more than single JSON");
	}


	return nil;

}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err:= json.Marshal(data);

	if err != nil {
		return err
	}

	if len(headers) > 0{
		for k, v:=range headers[0] {
			w.Header()[k] = v;
		}
	} 

	w.Header().Set("Content-Type","application/json");
	w.WriteHeader(status);
	
	_, err = w.Write(out);

	if err != nil {
		return err;
	}

	return nil;
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest;

	if len(status) > 0 {
		statusCode = status[0];
	}

	var payload JSONResponse;

	payload.Error = true;
	payload.Message = err.Error();
	

	return app.writeJSON(w, statusCode, payload);
}
