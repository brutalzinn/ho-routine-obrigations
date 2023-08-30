package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/brutalzinn/ho-routine-obrigations/firebase"
	"github.com/brutalzinn/ho-routine-obrigations/middlewares"
	"github.com/brutalzinn/ho-routine-obrigations/obrigation"
	"github.com/brutalzinn/ho-routine-obrigations/util"
	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

type ObrigationStartRequest struct {
	Id     string `json:"id"`
	Device string `json:"firebase_token"`
}

type ObrigationQRCodeRequest struct {
	Value  string `json:"value"`
	Device string `json:"firebase_token"`
}

type ObrigationQueuePending struct {
	Id     string
	Value  string
	Device string `json:"firebase_token"`
}

func main() {
	var obrigationPending *obrigation.Obrigation
	obrigationsQueue := make(chan ObrigationQueuePending)
	app := fiber.New()
	obrigation.Connect()
	obrigations, err := obrigation.ReadObrigations()
	if err != nil {
		log.Println("failed to read obrigations..")
		return
	}
	///HOME ASSISTANT
	haRoutes := app.Group("/ho", middlewares.ApiKeyMiddleware())
	{
		haRoutes.Get("/", websocket.New(func(c *websocket.Conn) {
			for {
				log.Println("WAITING CONFIRMATION ...")
				go func() {
					for item := range obrigationsQueue {
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
						obrigationPending = nil
						log.Println("Correct QR CODE")
					}
				}()
				time.Sleep(5 * time.Second)
			}
		}))

		haObrigations := haRoutes.Group("/obrigations")
		{
			haObrigations.Post("/", func(c *fiber.Ctx) {
				requestBody := new(ObrigationStartRequest)
				if err := c.BodyParser(requestBody); err != nil {
					c.SendStatus(400)
					return
				}
				found := false
				for _, obrigation := range obrigations {
					if obrigation.Id == requestBody.Id {
						obrigationPending = &obrigation
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
						"SOME OBRIGATION AT ROUTINE NEEDS YOUR ATTENTION "+obrigationPending.Name,
						"Tap this notification when you are ready to scan the QR CODE.")
					notify.Send()
				}()
				log.Println("Obrigation pending set to", obrigationPending)
				c.SendStatus(200)
			})

			haObrigations.Get("/", func(c *fiber.Ctx) {
				obrigations, err := obrigation.ReadObrigations()
				if err != nil {
					c.SendStatus(fiber.StatusNoContent)
					return
				}
				c.JSON(obrigations)
			})
		}

	}

	//mobile routes
	mobileRoutes := app.Group("/mobile", middlewares.ApiKeyMiddleware())
	{
		mobileObrigations := mobileRoutes.Group("/obrigation")
		{
			mobileObrigations.Post("/", func(c *fiber.Ctx) {
				requestBody := new(ObrigationQRCodeRequest)
				if err := c.BodyParser(requestBody); err != nil {
					c.SendStatus(fiber.StatusBadRequest)
					return
				}
				if obrigationPending == nil {
					c.JSON(fiber.Map{
						"confirmated": false,
						"message":     "No any obrigations pending at this moment.",
					})
					return
				}
				if obrigationPending.QrCode != requestBody.Value {
					c.JSON(fiber.Map{
						"confirmated": false,
						"message":     "Obrigation not foundß",
					})
					return
				}
				queue := ObrigationQueuePending{
					Value:  requestBody.Value,
					Device: requestBody.Device,
				}
				log.Println("Added obrigation to queue")
				go func() {
					obrigationsQueue <- queue
				}()
				c.Status(fiber.StatusCreated)
				c.JSON(fiber.Map{
					"confirmated": true,
					"message":     "Added to HO processing feedback.",
				})
			})

			mobileObrigations.Get("/", func(c *fiber.Ctx) {
				if obrigationPending != nil {
					c.JSON(obrigationPending)
					return
				}
				c.SendStatus(fiber.StatusNoContent)
			})
		}

	}
	///GET THE QR CODES TO PRINT
	app.Get("/qrcode", func(c *fiber.Ctx) {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		var sb strings.Builder
		for _, obrigation := range obrigations {
			sb.WriteString(fmt.Sprintf("<h1 style='margin:10px;'>%s</h1></br>", obrigation.Name))
			qrCode := util.GenerateQRCodeHtmlImageTag(obrigation.QrCode, 256)
			sb.WriteString(qrCode)
		}
		resultHtml := sb.String()
		c.SendString(resultHtml)
	})

	app.Get("/healthcheck", func(c *fiber.Ctx) {
		log.Println("This is just a healthcheck :)")
		c.SendStatus(fiber.StatusOK)
	})

	log.Fatal(app.Listen(":3030"))
}
