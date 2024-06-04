package main

import (
	"auth-svc/cmd/api/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPORT = ":80";

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
		Addr: webPORT,
		Handler: app.routes(),
	}

	log.Printf("Running Server on PORT:%s\n", webPORT);
	err := srv.ListenAndServe();

	if err != nil {
		log.Panic(err);
	}
	
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
	count := 0;

	for {
		connection, err := openDB(dsn);

		if err != nil {
			log.Println("Postgres not ready yet");
			count++;
		}	else {
			log.Println("Successfully connected to Postgres");
			return connection;
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