package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)


type Config struct {
	Mailer Mail
}

const webPORT = "80"


func main() {
	mailer := createMail();
	app := Config{
		Mailer: mailer,
	}
	
	log.Printf("Running mail service on PORT:%s",webPORT)
	
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err);
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain: os.Getenv("MAIL_DOMAIN"),
		Host: os.Getenv("MAIL_HOST"),
		Port: port,
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromName: os.Getenv("MAIL_FROM_NAME"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
	}

	return m;
}