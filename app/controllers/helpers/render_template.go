package helpers

import (
	"fmt"
	"html/template"
	"strings"
)

func RenderTemplate(templateStr string, data interface{}) (string, error) {
	// Parse the template from the string
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Create a buffer to store the rendered template
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	// Return the rendered HTML as a string
	return result.String(), nil
}
