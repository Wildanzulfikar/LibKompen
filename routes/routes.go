package routes

import (
	"LibKompen/controllers"
	"LibKompen/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Bebas Pustaka
	app.Get("/api/bebas-pustaka/:kode_user/history", middleware.Protected(), controllers.GetBebasPustakaHistory)
	app.Post("/api/bebas-pustaka/update", middleware.Protected(), controllers.UpdateBebasPustaka)

	// Mahasiswa Bebas Pustaka
	// Endpoint filter dan search via query parameter, contoh:
	// ?jurusan=07
	// ?status_pustaka=Bebas%20Pustaka
	// ?status_pinjaman=Lunas
	// ?tahun=2024
	app.Get("/api/mahasiswa-bebas-pustaka/", middleware.Protected(), controllers.GetMahasiswaBebasPustaka)

	// Loan
	app.Get("/api/loan/", controllers.GetAllLoan)
	app.Get("/api/loan/:loan_id", middleware.Protected(), controllers.GetLoanDetail)
	app.Delete("/api/loan/:loan_id", middleware.Protected(), controllers.DeleteLoanById)

	// Users
	app.Get("/api/users/", middleware.Protected(), controllers.UsersBebasPustakaGetAll)
	app.Delete("/api/users/:id", middleware.Protected(), controllers.DeleteUser)

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
