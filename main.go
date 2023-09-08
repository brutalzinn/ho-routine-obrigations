package main

import (
	"fmt"
	"log"
	"strings"

	database "github.com/brutalzinn/ho-routine-obrigations/db"
	"github.com/brutalzinn/ho-routine-obrigations/db/obrigation"
	"github.com/brutalzinn/ho-routine-obrigations/middlewares"
	"github.com/brutalzinn/ho-routine-obrigations/routes/homeassistant"
	"github.com/brutalzinn/ho-routine-obrigations/routes/mobile"
	"github.com/brutalzinn/ho-routine-obrigations/util"
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()
	database.Connect()
	obrigations, err := obrigation.GetObrigations()
	if err != nil {
		log.Println("failed to read obrigations..")
		return
	}
	///HOME ASSISTANT
	haRoutes := app.Group("/ho", middlewares.ApiKeyMiddleware())
	{
		haRoutes.Get("/", homeassistant.WebSocketHandler())
		haObrigations := haRoutes.Group("/obrigations")
		{
			haObrigations.Post("/start", homeassistant.StartObrigation)
			haObrigations.Get("/pending", homeassistant.GetPendingObrigation)
		}
	}

	//mobile routes
	mobileRoutes := app.Group("/mobile", middlewares.ApiKeyMiddleware())
	{
		mobileRoutes.Post("/register", mobile.InsertDevice)
		mobileObrigations := mobileRoutes.Group("/obrigation")
		{
			mobileObrigations.Post("/confirm", mobile.ConfirmObrigation)
			mobileObrigations.Get("/", mobile.GetPendingObrigation)
		}
	}

	app.Get("/", func(c *fiber.Ctx) {
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
		///this is to slowless my anxiety for half life 2 with RTX.
		stillAliveVideo := `<iframe width="560" height="315" src="https://www.youtube.com/embed/Y6ljFaKRTrI?si=xuARCBZ3ydxmsFwS" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>`
		log.Println("yeeep. I am alive!")
		c.SendString(stillAliveVideo)
	})

	log.Fatal(app.Listen(":3030"))
}
