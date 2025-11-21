package controllers

import (
    "github.com/gofiber/fiber/v2"
    "LibKompen/services"
)

func UpdateBebasPustaka(c *fiber.Ctx) error {
    return services.UpdateBebasPustakaService(c)
}
