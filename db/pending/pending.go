package pending

import (
	"fmt"
	"time"

	database "github.com/brutalzinn/ho-routine-obrigations/db"
)

type Pending struct {
	Id           int       `db:"id"`
	Waiting      bool      `db:"waiting"`
	IdDevice     int       `db:"id_device"`
	IdObrigation int       `db:"id_obrigation"`
	ExpireAt     time.Time `db:"expire_at"`
	CreateAt     time.Time `db:"create_at"`
	UpdateAt     time.Time `db:"update_at"`
}

func GetPendingsByDevice(idDevice int) (Pending, error) {
	var pending Pending
	rows, err := database.Conn.Query(`SELECT 
  id, 
  waiting, 
  id_device, 
  id_obrigation, 
  expire_at 
FROM 
  Pending 
WHERE 
  id_device = ?`, idDevice)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&pending.Id, &pending.Waiting, &pending.IdDevice, &pending.IdObrigation, &pending.ExpireAt)
		if err != nil {
			panic(err)
		}
		return pending, err
	}
	return pending, err
}

func GetPendings() ([]Pending, error) {
	var pendings []Pending
	rows, err := database.Conn.Query("SELECT id, waiting, id_device, id_obrigation, expire_at FROM Pending WHERE expire_at < NOW()")
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		var pending Pending
		err = rows.Scan(&pending.Id, &pending.Waiting, &pending.IdDevice, &pending.IdObrigation, &pending.ExpireAt)
		if err != nil {
			panic(err)
		}
		pendings = append(pendings, pending)
	}
	return pendings, err
}

func InsertPending(pending Pending) (bool, error) {
	rows, err := database.Conn.Exec(`INSERT INTO Pending(
  waiting, id_device, id_obrigation, 
  expire_at
) 
values 
  (?, ?, ?, ?)`, pending.Waiting, pending.IdDevice, pending.IdObrigation, pending.ExpireAt)

	changes, err := rows.RowsAffected()
	if changes == 0 {
		return false, err
	}
	return true, err
}

func UpdatePending(pending Pending) (bool, error) {
	rows, err := database.Conn.Exec(`UPDATE Pending SET waiting=?, expire_at=?, WHERE id=?`, pending.Waiting, pending.ExpireAt, pending.Id)
	changes, err := rows.RowsAffected()
	if changes == 0 {
		return false, err
	}
	return true, err
}
