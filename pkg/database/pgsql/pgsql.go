package pgsql

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to a DB: ", err)
	}

	// Actual connection check
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	// Close db connection on program termination
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit
		DB.Close()
	}()
}
