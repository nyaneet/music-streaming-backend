package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	ErrNoMatch    = fmt.Errorf("No matching record.")
	ErrNotAllowed = fmt.Errorf("Resource not allowed.")
)

type Database struct {
	Conn *sql.DB
}

func Initialize(host, username, password, database string, port int) (Database, error) {
	db := Database{}

	dbSource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}
