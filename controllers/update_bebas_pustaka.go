package controllers

import (
	"LibKompen/services"
	"github.com/gofiber/fiber/v2"
)

func UpdateBebasPustaka(c *fiber.Ctx) error {
	return services.UpdateBebasPustakaService(c)
}
