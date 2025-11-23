package main

import (
	"LibKompen/database"
	"LibKompen/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load .env file")
	}

	database.ConnectDB()

	app := fiber.New()

	app.Use(cors.New())

	app.Options("*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
