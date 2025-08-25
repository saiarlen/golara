package routes

import "github.com/test/myapp/framework"

// RegisterRoutes registers all application routes
func RegisterRoutes(app *framework.Golara) {
	// Register web routes (HTML pages)
	RegisterWebRoutes(app)
	
	// Register API routes (JSON responses)
	RegisterAPIRoutes(app)
}