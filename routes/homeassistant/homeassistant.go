package homeassistant

import (
	"log"
	"time"

	"github.com/gofiber/websocket"

	"github.com/brutalzinn/ho-routine-obrigations/db/device"
	"github.com/brutalzinn/ho-routine-obrigations/db/obrigation"
	"github.com/brutalzinn/ho-routine-obrigations/db/pending"
	"github.com/brutalzinn/ho-routine-obrigations/firebase"
	webmodels "github.com/brutalzinn/ho-routine-obrigations/models/web"
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
					log.Println("Verify if QR code is correct...", item.QrCodeValue)
					found := false
					for _, obrigation := range obrigations {
						if obrigation.QrCode == item.QrCodeValue {
							found = true
							break
						}
					}
					if !found {
						c.WriteJSON(fiber.Map{
							"confirmed": false,
							"message":   "incorrect QR CODE",
							"qr_code":   item.QrCodeValue,
						})
						log.Println("Wrong QR CODE")
						return
					}
					c.WriteJSON(fiber.Map{
						"confirmed": true,
						"message":   "OK",
						"qr_code":   item.QrCodeValue,
					})
					notify := firebase.New(
						item.TokenFirebase,
						"CONFIRMATION APPROVED",
						"You can close this application for now..")
					notify.Send()

					pendingObrigation := pending.Pending{
						IdObrigation: item.IdObrigation,
						ExpireAt:     time.Now(),
						Waiting:      false,
					}
					_, err := pending.UpdatePending(pendingObrigation)
					if err != nil {
						notify := firebase.New(
							item.TokenFirebase,
							"I THINK WE GOTTA A PROBLEM.",
							"THE OBRIGATION PENDING GOES WRONG AT CONFIRMATION STEP.")
						notify.Send()
						log.Println("failed to read obrigations..")
						return
					}
					log.Println("Correct QR CODE")
				}
			}()
			time.Sleep(5 * time.Second)
		}
	})
}

func StartObrigation(c *fiber.Ctx) {
	requestBody := new(webmodels.ObrigationStartRequest)
	if err := c.BodyParser(requestBody); err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "request wrong",
		})
	}
	///we need to handle error cases where we cant found any device
	mobileDevice, err := device.GetDevice(requestBody.DeviceName)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "device not found",
		})
		return
	}

	obrigationFound, err := obrigation.GetObrigation(requestBody.IdObrigation)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "obrigation not found",
		})
		return
	}
	pendent := pending.Pending{
		Waiting:      false,
		IdObrigation: obrigationFound.Id,
		IdDevice:     mobileDevice.Id,
	}
	pending.InsertPending(pendent)

	go func() {
		notify := firebase.New(
			mobileDevice.TokenFirebase,
			obrigationFound.Name,
			"Tap this notification when you are ready to scan the QR CODE.")
		notify.Send()
	}()
	log.Println("Obrigation pending set to", obrigationFound.Name)
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
