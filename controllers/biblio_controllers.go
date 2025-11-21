package controllers

import (
	"LibKompen/models"
	"LibKompen/services"
	"github.com/gofiber/fiber/v2"
)

func GetBiblio(c *fiber.Ctx) error {
	       biblios, err := services.GetAllBiblio()
	       if err != nil {
		       return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data biblio", "error": err.Error()})
	       }
	       return c.JSON(fiber.Map{"status": "success", "data": biblios})
}

func CreateBiblio(c *fiber.Ctx) error {
	       var biblio models.OpacBiblio
	       if err := c.BodyParser(&biblio); err != nil {
		       return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Gagal parse data"})
	       }
	       err := services.CreateBiblio(&biblio)
	       if err != nil {
		       return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data biblio", "error": err.Error()})
	       }
	       return c.JSON(fiber.Map{"status": "success", "data": biblio})
}
