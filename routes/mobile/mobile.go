package mobile

import (
	"log"

	"github.com/brutalzinn/ho-routine-obrigations/models"
	webmodels "github.com/brutalzinn/ho-routine-obrigations/models/web"
	"github.com/brutalzinn/ho-routine-obrigations/queue"
	"github.com/gofiber/fiber"
)

func GetObrigation(c *fiber.Ctx) {
	if queue.ObrigationPending != nil {
		c.JSON(queue.ObrigationPending)
		return
	}
	c.SendStatus(fiber.StatusNoContent)
}
func ConfirmObrigation(c *fiber.Ctx) {
	requestBody := new(webmodels.ObrigationConfirmRequest)
	if err := c.BodyParser(requestBody); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return
	}
	if queue.ObrigationPending == nil {
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "No any obrigations pending at this moment.",
		})
		return
	}
	if queue.ObrigationPending.QrCode != requestBody.Value {
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "Obrigation not foundß",
		})
		return
	}
	obrigationPending := models.ObrigationQueuePending{
		Value:  requestBody.Value,
		Device: requestBody.Device,
	}
	log.Println("Added obrigation to queue")
	go func() {
		queue.ObrigationsQueue <- obrigationPending
	}()
	c.Status(fiber.StatusCreated)
	c.JSON(fiber.Map{
		"confirmated": true,
		"message":     "Added to HO processing feedback.",
	})
}
