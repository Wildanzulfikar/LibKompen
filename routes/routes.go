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

	// Member routes
	app.Get("/members", controllers.GetAllMembers)
	app.Get("/members/:id", controllers.GetMemberByID)
	app.Post("/members", controllers.CreateMember)
	app.Put("/members/:id", controllers.UpdateMember)
	app.Delete("/members/:id", controllers.DeleteMember)

	// Loan routes
	app.Get("/loans", controllers.GetAllLoans)
	app.Get("/loans/:id", controllers.GetLoanByID)
	app.Get("/loans/member/:member_id", controllers.GetLoansByMember)
	app.Get("/loans/active", controllers.GetActiveLoans)
	app.Post("/loans", controllers.CreateLoan)
	app.Put("/loans/:id", controllers.UpdateLoan)
	app.Put("/loans/:id/return", controllers.ReturnLoan)
	app.Delete("/loans/:id", controllers.DeleteLoan)
}
