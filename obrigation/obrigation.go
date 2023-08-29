package obrigation

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Obrigation struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	QrCode    string    `db:"qr_code"`
	Mandatory bool      `db:"mandatory"`
	CreateAt  time.Time `db:"create_at"`
	UpdateAt  time.Time `db:"update_at"`
}

func Connect() error {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DB_CONFIG"))
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return err
}

func ReadObrigations() ([]Obrigation, error) {
	var obrigations []Obrigation
	rows, err := db.Query("SELECT id, name, mandatory, qr_code FROM Obrigations")
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		var obrigation Obrigation
		err = rows.Scan(&obrigation.Id, &obrigation.Name, &obrigation.Mandatory, &obrigation.QrCode)
		if err != nil {
			panic(err)
		}
		obrigations = append(obrigations, obrigation)
	}
	return obrigations, err

}
