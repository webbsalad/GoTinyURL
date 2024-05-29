package handler

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/skip2/go-qrcode"
	"github.com/webbsalad/GoTinyURL/config"
	"github.com/webbsalad/GoTinyURL/db"
	"github.com/webbsalad/GoTinyURL/db/operations"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.String()
	createApp().ServeHTTP(w, r)
}

func createApp() http.HandlerFunc {
	cfgDB, err := config.LoadConfig()
	if err != nil {
		log.Printf("Ошибка при чтении переменных окружения: %v\n", err)
		return nil
	}

	database := db.DBConnection{Config: cfgDB}

	if err := database.Connect(); err != nil {
		log.Printf("Ошибка при подключении к PostgreSQL: %v\n", err)
		return nil
	}

	app := fiber.New()

	app.Static("/", "var/templates")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("var/templates/index.html")
	})

	app.Post("/shorten", func(c *fiber.Ctx) error {
		var request struct {
			URL string `json:"url"`
		}
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse JSON",
			})
		}

		shortenedURL, err := operations.AddItem(&database, "urls", request.URL)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to shorten URL",
			})
		}

		qrCode, err := qrcode.Encode(shortenedURL, qrcode.Medium, 256)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate QR code",
			})
		}

		base64QRCode := base64.StdEncoding.EncodeToString(qrCode)
		qrCodeURL := "data:image/png;base64," + base64QRCode

		response := struct {
			ShortenedURL string `json:"shortenedUrl"`
			QRCodeURL    string `json:"qrCodeUrl"`
		}{
			ShortenedURL: shortenedURL,
			QRCodeURL:    qrCodeURL,
		}

		return c.JSON(response)
	})

	app.Get("/:shortenedKey", func(c *fiber.Ctx) error {
		shortenedKey := c.Params("shortenedKey")

		jsonData, err := operations.FetchKeyByValue(&database, "urls", shortenedKey)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("URL not found")
		}
		var result map[string]string
		if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse JSON")
		}

		originalURL, ok := result["key"]
		if !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("URL key not found in JSON")
		}

		if originalURL == c.Hostname() {
			return c.Status(fiber.StatusBadRequest).SendString("Cannot redirect to the same URL")
		}

		return c.Redirect(originalURL)
	})

	return adaptor.FiberApp(app)
}
