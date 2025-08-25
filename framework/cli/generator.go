package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Generator handles code generation for Golara framework
type Generator struct {
	basePath   string
	moduleName string
}

// NewGenerator creates a new generator instance
func NewGenerator(basePath string) *Generator {
	moduleName := getModuleName()
	return &Generator{
		basePath:   basePath,
		moduleName: moduleName,
	}
}

// getModuleName reads module name from go.mod
func getModuleName() string {
	file, err := os.Open("go.mod")
	if err != nil {
		return "myapp"
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module ")
		}
	}
	return "myapp"
}

// GenerateController creates a new controller
func (g *Generator) GenerateController(name string, webController ...bool) error {
	isWeb := len(webController) > 0 && webController[0]
	
	if isWeb {
		return g.generateWebController(name)
	}
	return g.generateAPIController(name)
}

// generateAPIController creates API controller
func (g *Generator) generateAPIController(name string) error {
	controllerTemplate := `package controllers

import (
	"{{.ModuleName}}/app/models"
	"{{.ModuleName}}/framework/database"
	"{{.ModuleName}}/framework/validation"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// {{.Name}}Controller handles {{.LowerName}} related requests (Laravel-style Resource Controller)
type {{.Name}}Controller struct {
	DB *database.DatabaseManager
}

// Index handles GET /{{.LowerName}} - Display a listing of the resource
func (ctrl *{{.Name}}Controller) Index(c *fiber.Ctx) error {
	var {{.LowerName}}s []models.{{.Name}}
	
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "15"))
	
	// Use query builder for pagination
	qb := database.NewQueryBuilder(ctrl.DB.Connection())
	result, err := qb.Table("{{.LowerName}}s").
		OrderBy("created_at", "DESC").
		Paginate(page, perPage, &{{.LowerName}}s)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch {{.LowerName}}s",
			"error":   err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data":    {{.LowerName}}s,
		"meta":    result,
	})
}

// Show handles GET /{{.LowerName}}/:id - Display the specified resource
func (ctrl *{{.Name}}Controller) Show(c *fiber.Ctx) error {
	id := c.Params("id")
	var {{.LowerName}} models.{{.Name}}
	
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Find(&{{.LowerName}}, id); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "{{.Name}} not found",
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data":    {{.LowerName}},
	})
}

// Store handles POST /{{.LowerName}} - Store a newly created resource
func (ctrl *{{.Name}}Controller) Store(c *fiber.Ctx) error {
	// Validate request
	validator := validation.NewValidator()
	// Add your validation rules here
	// validator.AddRule("name", validation.Required())
	
	if !validator.Validate(c) {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  validator.GetErrors(),
		})
	}
	
	// Create new {{.LowerName}}
	{{.LowerName}} := models.{{.Name}}{}
	if err := c.BodyParser(&{{.LowerName}}); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
		})
	}
	
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Create(&{{.LowerName}}); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create {{.LowerName}}",
			"error":   err.Error(),
		})
	}
	
	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "{{.Name}} created successfully",
		"data":    {{.LowerName}},
	})
}

// Update handles PUT /{{.LowerName}}/:id - Update the specified resource
func (ctrl *{{.Name}}Controller) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var {{.LowerName}} models.{{.Name}}
	
	// Find existing {{.LowerName}}
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Find(&{{.LowerName}}, id); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "{{.Name}} not found",
		})
	}
	
	// Parse update data
	var updateData models.{{.Name}}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
		})
	}
	
	// Update {{.LowerName}}
	if err := model.Update(&{{.LowerName}}, updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update {{.LowerName}}",
			"error":   err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "{{.Name}} updated successfully",
		"data":    {{.LowerName}},
	})
}

// Destroy handles DELETE /{{.LowerName}}/:id - Remove the specified resource
func (ctrl *{{.Name}}Controller) Destroy(c *fiber.Ctx) error {
	id := c.Params("id")
	var {{.LowerName}} models.{{.Name}}
	
	// Find existing {{.LowerName}}
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Find(&{{.LowerName}}, id); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "{{.Name}} not found",
		})
	}
	
	// Delete {{.LowerName}}
	if err := model.Delete(&{{.LowerName}}); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete {{.LowerName}}",
			"error":   err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "{{.Name}} deleted successfully",
	})
}
`

	data := struct {
		Name       string
		LowerName  string
		ModuleName string
	}{
		Name:       strings.Title(name),
		LowerName:  strings.ToLower(name),
		ModuleName: g.moduleName,
	}

	return g.generateFromTemplate(controllerTemplate, fmt.Sprintf("app/controllers/%s_controller.go", strings.ToLower(name)), data)
}

// generateWebController creates web controller with views
func (g *Generator) generateWebController(name string) error {
	webControllerTemplate := `package controllers

import (
	"{{.ModuleName}}/app/models"
	"{{.ModuleName}}/framework/database"
	"{{.ModuleName}}/framework/view"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// {{.Name}}Controller handles {{.LowerName}} web requests
type {{.Name}}Controller struct {
	DB *database.DatabaseManager
}

// Index displays listing page
func (ctrl *{{.Name}}Controller) Index(c *fiber.Ctx) error {
	var {{.LowerName}}s []models.{{.Name}}
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage := 15
	
	qb := database.NewQueryBuilder(ctrl.DB.Connection())
	result, err := qb.Table("{{.LowerName}}s").
		OrderBy("created_at", "DESC").
		Paginate(page, perPage, &{{.LowerName}}s)
	
	if err != nil {
		return c.Status(500).SendString("Error loading {{.LowerName}}s")
	}
	
	return view.View(c, "{{.LowerName}}/index", view.ViewData{
		"{{.LowerName}}s": {{.LowerName}}s,
		"pagination": result,
		"title": "{{.Name}}s",
	}, "layouts/app")
}

// Show displays single item
func (ctrl *{{.Name}}Controller) Show(c *fiber.Ctx) error {
	id := c.Params("id")
	var {{.LowerName}} models.{{.Name}}
	
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Find(&{{.LowerName}}, id); err != nil {
		return c.Status(404).SendString("{{.Name}} not found")
	}
	
	return view.View(c, "{{.LowerName}}/show", view.ViewData{
		"{{.LowerName}}": {{.LowerName}},
		"title": "{{.Name}} Details",
	}, "layouts/app")
}

// Create shows create form
func (ctrl *{{.Name}}Controller) Create(c *fiber.Ctx) error {
	return view.View(c, "{{.LowerName}}/create", view.ViewData{
		"title": "Create {{.Name}}",
	}, "layouts/app")
}

// Store handles form submission
func (ctrl *{{.Name}}Controller) Store(c *fiber.Ctx) error {
	{{.LowerName}} := models.{{.Name}}{}
	if err := c.BodyParser(&{{.LowerName}}); err != nil {
		return c.Redirect("/{{.LowerName}}s/create")
	}
	
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Create(&{{.LowerName}}); err != nil {
		return c.Redirect("/{{.LowerName}}s/create")
	}
	
	return c.Redirect("/{{.LowerName}}s")
}

// Edit shows edit form
func (ctrl *{{.Name}}Controller) Edit(c *fiber.Ctx) error {
	id := c.Params("id")
	var {{.LowerName}} models.{{.Name}}
	
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Find(&{{.LowerName}}, id); err != nil {
		return c.Status(404).SendString("{{.Name}} not found")
	}
	
	return view.View(c, "{{.LowerName}}/edit", view.ViewData{
		"{{.LowerName}}": {{.LowerName}},
		"title": "Edit {{.Name}}",
	}, "layouts/app")
}

// Update handles form update
func (ctrl *{{.Name}}Controller) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var {{.LowerName}} models.{{.Name}}
	
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Find(&{{.LowerName}}, id); err != nil {
		return c.Status(404).SendString("{{.Name}} not found")
	}
	
	var updateData models.{{.Name}}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Redirect("/{{.LowerName}}s/" + id + "/edit")
	}
	
	if err := model.Update(&{{.LowerName}}, updateData); err != nil {
		return c.Redirect("/{{.LowerName}}s/" + id + "/edit")
	}
	
	return c.Redirect("/{{.LowerName}}s")
}

// Destroy handles deletion
func (ctrl *{{.Name}}Controller) Destroy(c *fiber.Ctx) error {
	id := c.Params("id")
	var {{.LowerName}} models.{{.Name}}
	
	model := database.NewModel(ctrl.DB.Connection())
	if err := model.Find(&{{.LowerName}}, id); err != nil {
		return c.Status(404).SendString("{{.Name}} not found")
	}
	
	if err := model.Delete(&{{.LowerName}}); err != nil {
		return c.Redirect("/{{.LowerName}}s")
	}
	
	return c.Redirect("/{{.LowerName}}s")
}
`

	data := struct {
		Name       string
		LowerName  string
		ModuleName string
	}{
		Name:       strings.Title(name),
		LowerName:  strings.ToLower(name),
		ModuleName: g.moduleName,
	}

	return g.generateFromTemplate(webControllerTemplate, fmt.Sprintf("app/controllers/%s_controller.go", strings.ToLower(name)), data)
}

// GenerateMigration creates a new migration file
func (g *Generator) GenerateMigration(name string) error {
	timestamp := time.Now().Format("2006_01_02_150405")
	fileName := fmt.Sprintf("%s_create_%s_table", timestamp, strings.ToLower(name))
	
	migrationTemplate := `package migrations

import (
	"gorm.io/gorm"
)

// {{.Name}}Migration creates {{.LowerName}} table
func {{.Name}}Migration() (string, func(*gorm.DB) error, func(*gorm.DB) error) {
	return "{{.FileName}}",
		// Up
		func(db *gorm.DB) error {
			return db.Exec(` + "`" + `
				CREATE TABLE {{.LowerName}}s (
					id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
					created_at TIMESTAMP NULL,
					updated_at TIMESTAMP NULL,
					deleted_at TIMESTAMP NULL,
					INDEX idx_{{.LowerName}}s_deleted_at (deleted_at)
				)
			` + "`" + `).Error
		},
		// Down
		func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS {{.LowerName}}s").Error
		}
}
`

	data := struct {
		Name      string
		LowerName string
		FileName  string
	}{
		Name:      strings.Title(name),
		LowerName: strings.ToLower(name),
		FileName:  fileName,
	}

	return g.generateFromTemplate(migrationTemplate, fmt.Sprintf("database/migrations/%s.go", fileName), data)
}

// GenerateModel creates a new model file
func (g *Generator) GenerateModel(name string) error {
	modelTemplate := `package models

import (
	"gorm.io/gorm"
)

// {{.Name}} model
type {{.Name}} struct {
	gorm.Model
	// Add your fields here
}

// TableName returns the table name
func ({{.Name}}) TableName() string {
	return "{{.LowerName}}s"
}
`

	data := struct {
		Name      string
		LowerName string
	}{
		Name:      strings.Title(name),
		LowerName: strings.ToLower(name),
	}

	return g.generateFromTemplate(modelTemplate, fmt.Sprintf("app/models/%s.go", strings.ToLower(name)), data)
}

// GenerateMiddleware creates a new middleware file
func (g *Generator) GenerateMiddleware(name string) error {
	middlewareTemplate := `package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// {{.Name}} middleware
func {{.Name}}() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Add your middleware logic here
		return c.Next()
	}
}
`

	data := struct {
		Name string
	}{
		Name: strings.Title(name),
	}

	return g.generateFromTemplate(middlewareTemplate, fmt.Sprintf("app/middleware/%s.go", strings.ToLower(name)), data)
}

// GenerateJob creates a new job file
func (g *Generator) GenerateJob(name string) error {
	jobTemplate := `package jobs

import (
	"{{.ModuleName}}/framework/queue"
	"log"
)

// {{.Name}}Job handles {{.LowerName}} processing
type {{.Name}}Job struct {
	queue.BaseJob
	// Add your job fields here
}

func (j *{{.Name}}Job) Handle() error {
	log.Printf("Processing {{.Name}} job")
	// Add your job logic here
	return nil
}

func New{{.Name}}Job() *{{.Name}}Job {
	return &{{.Name}}Job{
		BaseJob: queue.BaseJob{
			Name: "{{.LowerName}}",
			Payload: make(map[string]interface{}),
		},
	}
}
`

	data := struct {
		Name       string
		LowerName  string
		ModuleName string
	}{
		Name:       strings.Title(name),
		LowerName:  strings.ToLower(name),
		ModuleName: g.moduleName,
	}

	return g.generateFromTemplate(jobTemplate, fmt.Sprintf("app/jobs/%s_job.go", strings.ToLower(name)), data)
}

// GenerateView creates a new view template
func (g *Generator) GenerateView(name string) error {
	viewTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    {{asset "css/app.css"}}
</head>
<body>
    <div class="container">
        <h1>{{.Title}}</h1>
        
        <!-- Your content here -->
        <p>Welcome to {{.Name}} view!</p>
        
        <!-- Example data display -->
        {{if .Data}}
        <div class="data">
            {{range .Data}}
            <div class="item">{{.}}</div>
            {{end}}
        </div>
        {{end}}
    </div>
    
    {{asset "js/app.js"}}
</body>
</html>
`

	layoutTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - {{.AppName}}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    {{csrf}}
    <link href="{{asset "css/app.css"}}" rel="stylesheet">
</head>
<body>
    <nav class="navbar">
        <div class="container">
            <a href="{{route "home"}}" class="brand">{{.AppName}}</a>
        </div>
    </nav>
    
    <main class="main">
        {{template "content" .}}
    </main>
    
    <script src="{{asset "js/app.js"}}"></script>
</body>
</html>
`

	data := struct {
		Name    string
		Title   string
		AppName string
	}{
		Name:    strings.Title(name),
		Title:   strings.Title(name),
		AppName: "Golara App",
	}

	// Create layout if it doesn't exist
	layoutPath := "resources/views/layouts/app.html"
	if _, err := os.Stat(layoutPath); os.IsNotExist(err) {
		if err := g.generateFromTemplate(layoutTemplate, layoutPath, data); err != nil {
			return err
		}
	}

	return g.generateFromTemplate(viewTemplate, fmt.Sprintf("resources/views/%s.html", strings.ToLower(name)), data)
}

func (g *Generator) generateFromTemplate(templateStr, filePath string, data interface{}) error {
	tmpl, err := template.New("generator").Parse(templateStr)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(g.basePath, filePath)
	dir := filepath.Dir(fullPath)
	
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}