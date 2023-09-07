package homeassistant

import (
	"log"
	"time"

	"github.com/gofiber/websocket"

	"github.com/brutalzinn/ho-routine-obrigations/firebase"
	webmodels "github.com/brutalzinn/ho-routine-obrigations/models/web"
	"github.com/brutalzinn/ho-routine-obrigations/obrigation"
	"github.com/brutalzinn/ho-routine-obrigations/queue"
	"github.com/gofiber/fiber"
)

func WebSocketHandler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		obrigations, err := obrigation.GetObrigations()
		if err != nil {
			log.Println("failed to read obrigations..")
			return
		}
		for {
			log.Println("WAITING CONFIRMATION ...")
			go func() {
				for item := range queue.ObrigationsQueue {
					log.Println("Received a confirmation request ...", item.Id)
					log.Println("Verify if QR code is correct...", item.Value)
					found := false
					for _, obrigation := range obrigations {
						if obrigation.QrCode == item.Value {
							found = true
							break
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
					notify := firebase.New(
						item.Device,
						"CONFIRMATION APPROVED",
						"You can close this application for now..")
					notify.Send()
					queue.ObrigationPending = nil
					log.Println("Correct QR CODE")
				}
			}()
			time.Sleep(5 * time.Second)
		}
	})
}

func StartObrigation(c *fiber.Ctx) {
	obrigations, err := obrigation.GetObrigations()
	if err != nil {
		log.Println("failed to read obrigations..")
		return
	}
	requestBody := new(webmodels.ObrigationStartRequest)
	if err := c.BodyParser(requestBody); err != nil {
		c.SendStatus(400)
		return
	}
	found := false
	for _, obrigation := range obrigations {
		if obrigation.Id == requestBody.Id {
			queue.ObrigationPending = &obrigation
			found = true
			break
		}
	}
	if !found {
		log.Println("Obrigation not found")
		c.SendStatus(400)
		return
	}
	go func() {
		notify := firebase.New(
			requestBody.Device,
			"SOME OBRIGATION AT ROUTINE NEEDS YOUR ATTENTION "+queue.ObrigationPending.Name,
			"Tap this notification when you are ready to scan the QR CODE.")
		notify.Send()
	}()
	log.Println("Obrigation pending set to", queue.ObrigationPending)
	c.SendStatus(200)
}

func GetPendingObrigation(c *fiber.Ctx) {
	obrigations, err := obrigation.GetObrigations()
	if err != nil {
		c.SendStatus(fiber.StatusNoContent)
		return
	}
	c.JSON(obrigations)
}
