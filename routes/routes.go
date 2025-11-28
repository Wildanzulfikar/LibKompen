package routes

import (
	"LibKompen/controllers"
	"LibKompen/middleware"

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
	tenggat.Post("/", middleware.Protected(), controllers.CreateTenggat)
	tenggat.Delete("/:id_tenggat_waktu", controllers.DeleteTenggat)
	tenggat.Put("/:id_tenggat_waktu", controllers.UpdateTenggat)
	// Bebas Pustaka
	app.Get("/api/bebas-pustaka/:kode_user/history", middleware.Protected(), controllers.GetBebasPustakaHistory)
	app.Post("/api/bebas-pustaka/update", middleware.Protected(), controllers.UpdateBebasPustaka)

	// Mahasiswa Bebas Pustaka
	app.Get("/api/mahasiswa-bebas-pustaka/", controllers.GetMahasiswaBebasPustaka)

	// Loan
	app.Get("/api/loan/", controllers.GetAllLoan) // new endpoint for all loans
	app.Get("/api/loan/:loan_id", middleware.Protected(), controllers.GetLoanDetail)
	app.Delete("/api/loan/:loan_id", middleware.Protected(), controllers.DeleteLoanById)

	// Users
	app.Get("/api/users/", middleware.Protected(), controllers.UsersBebasPustakaGetAll)
	app.Delete("/api/users/:id", middleware.Protected(), controllers.DeleteUsersBebasPustaka)

	// Auth
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/logout", middleware.Protected(), controllers.Logout)
	app.Get("/api/me", middleware.Protected(), controllers.Me)

	// Biblio & Approvals
	// app.Get("/api/biblio", controllers.GetBiblio)
	// app.Post("/api/biblio", controllers.CreateBiblio)
	// app.Get("/api/approvals", controllers.GetApprovals)
	// app.Post("/api/approvals", controllers.CreateApproval)

}
