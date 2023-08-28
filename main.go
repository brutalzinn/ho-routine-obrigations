package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

const (
	OBRIGATIONS_FILE = "obrigations.json"
)

type ObrigationStartRequest struct {
	Id string `json:"id"`
}

type ObrigationQRCodeRequest struct {
	Value string `json:"value"`
}

type Obrigation struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	QrCode    string `json:"qr_code"`
	Mandatory bool   `json:"mandatory"`
}
type ObrigationQueuePending struct {
	Id    string
	Value string
}

func ReadObrigations() ([]Obrigation, error) {
	jsonFile, err := os.Open(OBRIGATIONS_FILE)
	if err != nil {
		log.Printf("obrigations file not found")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var obrigations []Obrigation
	err = json.Unmarshal(byteValue, &obrigations)
	return obrigations, err

}
func main() {
	obrigations, err := ReadObrigations()
	var obrigationPending ObrigationQueuePending
	if err != nil {
		log.Println("failed to read obrigations..")
		return
	}
	obrigationsQueue := make(chan ObrigationQueuePending)
	app := fiber.New()

	app.Get("/:obrigationId", websocket.New(func(c *websocket.Conn) {
		obrigationId := c.Params("obrigationId")
		for {
			log.Println("WAITING CONFIRMATION for ...", obrigationId)
			go func() {
				for item := range obrigationsQueue {
					log.Println("Received a confirmation request ...", item.Id)
					log.Println("Verify if QR code is correct...", item.Value)
					found := false
					for _, obrigation := range obrigations {
						if obrigation.QrCode == item.Value {
							found = true
						}
					}
					if !found {
						c.WriteJSON(fiber.Map{
							"confirmed": false,
							"message":   "incorrect QR CODE",
							"qr_code":   item.Value,
						})
						log.Println("Wrong QR CODE")
						return
					}
					c.WriteJSON(fiber.Map{
						"confirmed": true,
						"message":   "OK",
						"qr_code":   item.Value,
					})
					log.Println("Correct QR CODE")
				}
			}()
			time.Sleep(5 * time.Second)
		}
	}))
	app.Post("/obrigation/confirm", func(c *fiber.Ctx) {
		requestBody := new(ObrigationQRCodeRequest)
		if err := c.BodyParser(requestBody); err != nil {
			c.SendStatus(400)
			return
		}
		if obrigationPending.Value != requestBody.Value {
			c.JSON(fiber.Map{
				"message": "no any obrigations with this qr code",
			})
			return
		}
		queue := ObrigationQueuePending{
			Value: requestBody.Value,
		}
		log.Println("Added obrigation to queue")
		go func() {
			obrigationsQueue <- queue
		}()
		c.SendStatus(200)
	})

	app.Post("/obrigation/start", func(c *fiber.Ctx) {
		requestBody := new(ObrigationStartRequest)
		if err := c.BodyParser(requestBody); err != nil {
			c.SendStatus(400)
			return
		}
		found := false
		for _, obrigation := range obrigations {
			if obrigation.Id == requestBody.Id {
				found = true
			}
		}
		if !found {
			log.Println("Obrigation not found")
			c.SendStatus(400)
			return
		}
		obrigationPending = ObrigationQueuePending{
			Id: requestBody.Id,
		}
		log.Println("Obrigation pending set to", obrigationPending)
		c.SendStatus(200)
	})

	app.Get("/healthcheck", func(c *fiber.Ctx) {
		log.Println("This is just a healthcheck :)")
		c.SendStatus(200)
	})

	log.Fatal(app.Listen(":3030"))
}
