package routes

import (
	"LibKompen/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/biblio", controllers.GetBiblio)
	app.Post("/biblio", controllers.CreateBiblio)

	app.Get("/approvals", controllers.GetApprovals)
	app.Post("/approvals", controllers.CreateApproval)
}
