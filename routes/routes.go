package routes

import (
	"LibKompen/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	users := app.Group("/users")

	users.Get("/", controllers.UsersBebasPustakaGetAll)
	users.Post("/", controllers.CreateUsersBebasPustaka)
	users.Delete("/:id_users", controllers.DeleteUsersBebasPustaka)

	tenggat := app.Group("/tenggat")

	tenggat.Get("/", controllers.GetAllTenggat)
	tenggat.Get("/aktif", controllers.GetActiveTenggat)
	tenggat.Post("/", controllers.CreateTenggat)
	tenggat.Delete("/:id_tenggat_waktu", controllers.DeleteTenggat)
	tenggat.Put("/:id_tenggat_waktu", controllers.UpdateTenggat)
}
