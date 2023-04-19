package main

import (
	"astro/persistence"
	"astro/vpn/api/rest"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"strings"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOriginsFunc: nil,
		AllowOrigins:     "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${ip}:${port} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	app.Get("/metrics", monitor.New(monitor.Config{APIOnly: true}))

	api := app.Group("/astro")

	persists, err := persistence.NewPersists()
	if err != nil {
		log.Fatal(err)
	}

	restVpn := rest.NewRestAPI(persists.DbVPN, persists.DBConfig, api.Group("/vpn"))
	restVpn.Handlers()

	log.Fatal(app.Listen(":80"))
}
