package main

import (
	"auth-svc/cmd/api/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPORT = "8080";

const MaxDBConnTries int = 15;
const RetryDBConnIn time.Duration = 5;

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main() {
	log.Println("Booting up authentication server...");	
	
	conn := connectToDB();

	if conn == nil {
		log.Panic("Database connection failed.");
	}

	app := Config{
		DB: conn,
		Models: data.New(conn),
	};
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",webPORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe();

	if err != nil {
		log.Panic(err);
	}

	log.Println("Auth-SVC running on PORT:%s", webPORT)
	
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn);

	if err != nil {
		return nil, err
	}

	err = db.Ping();

	if err != nil {
		return nil, err
	}

	return db, nil
}


func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN");
	var count int = 0;

	for {
		conn, err := openDB(dsn);

		if err != nil {
			log.Println("Postgres not ready yet");
			count++;
		}	else {
			log.Println("Successfully connected to Postgres");
			return conn;
		}

		if count > MaxDBConnTries {
			log.Println("Couldn't connec to DB")
			return nil;
		}

		log.Println("Retrying in 5 seconds");
		time.Sleep(RetryDBConnIn*time.Second);
		continue;
	}
}