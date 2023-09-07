package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var Conn *sql.DB

func Connect() error {
	var err error
	Conn, err = sql.Open("mysql", os.Getenv("DB_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}
	pingErr := Conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return err
}
