package routes

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes(app *fiber.App) {
	//api := app.Group("/api")

	///////////////////////////// public Api routes //////////////////////////////

	// Admin auth
	// api.Post("/admin/admin-auth", adminkycauth.Admin)
	// api.Post("/admin/admin-otpverify", adminkycauth.AdminOtpVerify)

	//privateAdmin := api.Group("/admin", middlewares.AdminTokenAuth)

	///////////////////////////////////////////  Private ADMIN      ///////////////////////////////////////////////////////////////////

	//privateAdmin.Get("/dashboard-analytics", admin.AdminDashboard)

}
