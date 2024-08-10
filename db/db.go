package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB(dsn string) error {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("sql.Open err: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("DB.Ping err: %v", err)
	}
	return nil
}
