package obrigation

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Obrigation struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	QrCode    string    `db:"qr_code"`
	Mandatory bool      `db:"mandatory"`
	CreateAt  time.Time `db:"create_at"`
	UpdateAt  time.Time `db:"update_at"`
}

func ReadObrigations() ([]Obrigation, error) {
	var obrigations []Obrigation
	db, err := sql.Open("mysql", os.Getenv("DB_CONFIG"))
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Obrigations")
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
