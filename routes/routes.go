package routes

import (
	"LibKompen/controllers"
	"LibKompen/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/users", controllers.UsersBebasPustakaGetAll)

	// Auth routes
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/logout", middleware.Protected(), controllers.Logout)
	app.Get("/api/me", middleware.Protected(), controllers.Me)

	// app.Get("/biblio", controllers.GetBiblio)
	// app.Post("/biblio", controllers.CreateBiblio)

	// app.Get("/approvals", controllers.GetApprovals)
	// app.Post("/approvals", controllers.CreateApproval)
}
