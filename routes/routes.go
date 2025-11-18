package routes

import (
	"LibKompen/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/users", controllers.UsersBebasPustakaGetAll)
	app.Post("/users", controllers.CreateUsersBebasPustaka)

	app.Delete("/users/:id_users", controllers.DeleteUsersBebasPustaka)
}
