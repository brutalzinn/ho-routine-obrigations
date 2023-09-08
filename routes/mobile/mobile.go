package mobile

import (
	"log"

	"github.com/brutalzinn/ho-routine-obrigations/db/device"
	"github.com/brutalzinn/ho-routine-obrigations/db/pending"
	"github.com/brutalzinn/ho-routine-obrigations/models"
	webmodels "github.com/brutalzinn/ho-routine-obrigations/models/web"
	"github.com/brutalzinn/ho-routine-obrigations/queue"
	"github.com/gofiber/fiber"
)

func GetPendingObrigation(c *fiber.Ctx) {
	pendings, err := pending.GetPendings()
	if err != nil {
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "No obrigations found at this moment",
		})
		return
	}
	c.JSON(pendings)
}

func InsertDevice(c *fiber.Ctx) {

	requestBody := new(webmodels.MobileRegisterDevice)
	if err := c.BodyParser(requestBody); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return
	}
	mobileDevice := device.Device{
		Name:          requestBody.Name,
		TokenFirebase: requestBody.FirebaseToken,
	}
	pendingObrigation, err := device.InsertDevice(mobileDevice)
	if err != nil {
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "Cant register this device now.",
		})
		return
	}
	c.JSON(pendingObrigation)

}
func ConfirmObrigation(c *fiber.Ctx) {
	requestBody := new(webmodels.MobileConfirmRequest)
	if err := c.BodyParser(requestBody); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return
	}

	mobileDevice, err := device.GetDevice(requestBody.Value)
	if err != nil {
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "Device not found",
		})
		return
	}
	pendingObrigation, err := pending.GetPendingsByDevice(mobileDevice.Id)
	if err != nil {
		c.JSON(fiber.Map{
			"confirmated": false,
			"message":     "No pending found at this moment",
		})
		return
	}
	obrigationPending := models.ObrigationQueuePending{
		QrCodeValue:   requestBody.Value,
		IdObrigation:  pendingObrigation.IdObrigation,
		IdDevice:      pendingObrigation.IdDevice,
		DeviceName:    mobileDevice.Name,
		TokenFirebase: mobileDevice.TokenFirebase,
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
