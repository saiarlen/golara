package controllers

import (
	"github.com/test/myapp/app/models"
	"github.com/test/myapp/framework/database"
	"github.com/test/myapp/framework/events"
	"github.com/test/myapp/framework/validation"

	"github.com/gofiber/fiber/v2"
)

// UserController handles user operations
type UserController struct {
	DB     *database.DatabaseManager
	Events *events.EventDispatcher
}

// Index returns paginated users
func (uc *UserController) Index(c *fiber.Ctx) error {
	var users []models.User

	qb := database.NewQueryBuilder(uc.DB.Connection())
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 15)

	result, err := qb.Table("users").
		Where("status", "=", "active").
		OrderBy("created_at", "DESC").
		Paginate(page, perPage, &users)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	return c.JSON(result)
}

// Store creates a new user
func (uc *UserController) Store(c *fiber.Ctx) error {
	validator := validation.NewValidator()

	validator.AddRule("name", validation.Required()).
		AddRule("name", validation.Min(2)).
		AddRule("email", validation.Required()).
		AddRule("email", validation.Email()).
		AddRule("password", validation.Required()).
		AddRule("password", validation.Min(8))

	if !validator.Validate(c) {
		return c.Status(400).JSON(fiber.Map{"errors": validator.GetErrors()})
	}

	data := validator.GetData()

	user := models.User{
		Name:     data["name"].(string),
		Email:    data["email"].(string),
		Password: data["password"].(string),
		Status:   "active",
	}

	model := database.NewModel(uc.DB.Connection())
	if err := model.Create(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Dispatch event
	event := events.NewUserRegisteredEvent(string(rune(user.ID)), user.Email)
	uc.Events.DispatchAsync(event)

	return c.Status(201).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    user,
	})
}
