package view

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// Engine implements Fiber's Views interface with Laravel Blade-like features
type Engine struct {
	directory string
	extension string
	funcmap   template.FuncMap
}

// New creates a new template engine
func New(directory, extension string) *Engine {
	engine := &Engine{
		directory: directory,
		extension: extension,
		funcmap:   make(template.FuncMap),
	}
	engine.addHelpers()
	return engine
}

// addHelpers adds Laravel-like helper functions
func (e *Engine) addHelpers() {
	e.funcmap["asset"] = func(path string) string {
		return "/assets/" + path
	}
	
	e.funcmap["route"] = func(name string) string {
		return "/" + name
	}
	
	e.funcmap["csrf"] = func() template.HTML {
		return template.HTML(`<input type="hidden" name="_token" value="csrf-token">`)
	}
}

// Load loads all templates
func (e *Engine) Load() error {
	return nil
}

// Render renders a template
func (e *Engine) Render(out io.Writer, name string, binding interface{}, layout ...string) error {
	tmplPath := filepath.Join(e.directory, name+e.extension)
	
	tmpl := template.New(name).Funcs(e.funcmap)
	
	if len(layout) > 0 {
		layoutPath := filepath.Join(e.directory, "layouts", layout[0]+e.extension)
		tmpl, _ = tmpl.ParseFiles(layoutPath, tmplPath)
	} else {
		tmpl, _ = tmpl.ParseFiles(tmplPath)
	}
	
	return tmpl.Execute(out, binding)
}

// ViewData holds data for views
type ViewData map[string]interface{}

// View helper for controllers
func View(c *fiber.Ctx, template string, data ViewData, layout ...string) error {
	if len(layout) > 0 {
		return c.Render(template, data, layout[0])
	}
	return c.Render(template, data)
}