package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EnsureDir creates directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists checks if file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ToSnakeCase converts string to snake_case
func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ToPascalCase converts string to PascalCase
func ToPascalCase(str string) string {
	parts := strings.Split(str, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}

// GenerateTimestamp generates timestamp for migrations
func GenerateTimestamp() string {
	return time.Now().Format("2006_01_02_150405")
}

// GetProjectRoot returns project root directory
func GetProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// Look for go.mod file
	for {
		if FileExists(filepath.Join(wd, "go.mod")) {
			return wd, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	
	return "", fmt.Errorf("project root not found")
}