package main

import (
	"fmt"
	"log"
	"net/http"
)


type Config struct {

}

const webPORT = "80"


func main() {
	app := Config{}
	
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