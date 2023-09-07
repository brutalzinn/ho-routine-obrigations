package obrigation

import (
	"fmt"
	"time"

	database "github.com/brutalzinn/ho-routine-obrigations/db"
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

func GetObrigations() ([]Obrigation, error) {
	var obrigations []Obrigation
	rows, err := database.Conn.Query("SELECT id, name, mandatory, qr_code FROM Obrigations")
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
