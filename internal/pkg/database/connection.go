package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func ConnectDB() *sql.DB {
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	dbname := os.Getenv("PGDBNAME")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("Can't connect to database with err: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database not responding with err: %s", err)
	}

	log.Println("Success Connect To Database")

	return db
}
