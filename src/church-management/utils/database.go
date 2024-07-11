package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	const (
		host     = "viaduct.proxy.rlwy.net"
		port     = 50913
		user     = "postgres"
		password = "DwVqApocFjDrFUqfYSbZhrDXwveKNFyq"
		dbname   = "railway"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")
	DB = db
	return db, nil
}
