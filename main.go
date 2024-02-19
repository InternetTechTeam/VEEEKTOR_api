package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func initDBConnection() {
	var err error
	db, err = sql.Open("pgx", os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatal("Failed to connect to a DB: ", err)
	}

	// Actual connection check
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
}

func getMultiplexer() *http.ServeMux {
	mux := http.NewServeMux()

	// Mux handlers

	return mux
}

func main() {
	log.Printf("VEEEKTOR_api is starting...")

	initDBConnection()
	defer db.Close()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := getMultiplexer()
	server := &http.Server{Addr: ":8080", Handler: mux}
	defer server.Close()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
