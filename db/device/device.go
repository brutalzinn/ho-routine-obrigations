package device

import (
	"fmt"
	"time"

	database "github.com/brutalzinn/ho-routine-obrigations/db"
)

type Device struct {
	Id            int       `db:"id"`
	Name          string    `db:"name"`
	TokenFirebase string    `db:"token_firebase"`
	CreateAt      time.Time `db:"create_at"`
	UpdateAt      time.Time `db:"update_at"`
}

func GetDevices() ([]Device, error) {
	var devices []Device
	rows, err := database.Conn.Query("SELECT id, name, token_firebase FROM devices")
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		var device Device
		err = rows.Scan(&device.Id, &device.Name, &device.TokenFirebase)
		if err != nil {
			panic(err)
		}
		devices = append(devices, device)
	}
	return devices, err
}

func GetDevice(name string) (Device, error) {
	var device Device
	rows, err := database.Conn.Query("SELECT id, name, token_firebase FROM Devices where name = ?", name)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&device.Id, &device.Name, &device.TokenFirebase)
		if err != nil {
			panic(err)
		}
		return device, err
	}
	return device, err
}

func InsertDevice(device Device) (bool, error) {
	rows, err := database.Conn.Exec("INSERT INTO Devices (name, token_firebase) VALUES (?,?)", device.Name, device.TokenFirebase)
	if err != nil {
		fmt.Print(err)
	}
	changes, err := rows.RowsAffected()
	if changes == 0 {
		return false, err
	}
	return true, err
}
