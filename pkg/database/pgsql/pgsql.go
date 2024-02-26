package pgsql

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func init() {
	// db url should be get by os.Getenv()
	var err error
	DB, err = sql.Open("pgx", "postgresql://veeektor:766180@localhost:5432/veeektor_db")
	if err != nil {
		log.Fatal("Failed to connect to a DB: ", err)
	}

	// Actual connection check
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
}
