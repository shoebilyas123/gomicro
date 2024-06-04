package main

import (
	"fmt"
	"log"
	"net/http"
)


type Config struct {}

const brokerAddr = "80"

func main(){
	app := Config{};

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",brokerAddr),
		Handler: app.routes(),
	}

	log.Printf("Server listening on PORT:%s\n",brokerAddr);
	err := srv.ListenAndServe();

	if err != nil {
		log.Fatal(err);
	}
	
}