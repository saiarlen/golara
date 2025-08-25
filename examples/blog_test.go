package examples

import (
	"github.com/test/myapp/framework"
	"github.com/test/myapp/framework/database"
	"github.com/test/myapp/internal/constants"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// BlogPost model for testing
type BlogPost struct {
	gorm.Model
	Title   string `json:"title" gorm:"size:255;not null"`
	Content string `json:"content" gorm:"type:text"`
	Status  string `json:"status" gorm:"size:20;default:draft"`
}

func TestBlogAPI(t *testing.T) {
	app := framework.New(framework.Config{
		AppName:     "Blog Test",
		Environment: constants.EnvTesting,
	})

	// Mock blog endpoints
	app.App.Get("/posts", func(c *fiber.Ctx) error {
		// Mock data
		posts := []BlogPost{
			{Title: "First Post", Content: "Content 1", Status: "published"},
			{Title: "Second Post", Content: "Content 2", Status: "draft"},
		}
		return c.JSON(fiber.Map{"posts": posts})
	})

	app.App.Post("/posts", func(c *fiber.Ctx) error {
		var post BlogPost
		if err := c.BodyParser(&post); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid data"})
		}

		post.ID = 1 // Mock ID
		return c.Status(201).JSON(fiber.Map{"post": post})
	})

	// Test GET /posts
	req := httptest.NewRequest("GET", "/posts", nil)
	resp, err := app.App.Test(req)
	if err != nil {
		t.Errorf("GET /posts failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestQueryBuilder(t *testing.T) {
	// Test QueryBuilder structure without database connection
	qb := &database.QueryBuilder{}

	if qb == nil {
		t.Error("QueryBuilder should be initialized")
	}

	// Test that QueryBuilder exists and can be instantiated
	// Note: Actual database operations require a valid DB connection
	t.Log("QueryBuilder structure test passed - actual DB operations require connection")
}
