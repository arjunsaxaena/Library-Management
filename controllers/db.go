package controllers

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB // db holds the connection to the PostgreSQL database

func InitDB() error {
	var err error
	connStr := "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"
	db, err = sqlx.Connect("postgres", connStr) // If success connection assigned to db
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return nil
}

func GetDB() *sqlx.DB {
	return db
}
