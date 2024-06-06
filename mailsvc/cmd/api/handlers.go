package main

import (
	"log"
	"net/http"
)

type MailMessage struct {
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	var reqPayload MailMessage

	err := app.readJSON(w, r, &reqPayload);

	if err != nil {
		log.Printf("CHECK1:%T\n",err);
		app.errorJSON(w, err);
		return;
	}

	msg := Message{
		From: reqPayload.From,
		To: reqPayload.To,
		Subject: reqPayload.Subject,
		Data: reqPayload.Message,
	}

	log.Printf("%v\n",reqPayload);
	
	err = app.Mailer.sendSMTPMessage(msg);

	if err != nil {
		log.Printf("CHECK2:%v\n",err);
		app.errorJSON(w, err);
		return;
	}

	var payload JSONResponse

	payload.Error = false;
	payload.Message = "Sent to"+reqPayload.To;

	app.writeJSON(w, http.StatusAccepted, payload)
}