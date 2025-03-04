package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Validator holds the validation rules and data
type Validator struct {
	Rules  map[string]string
	Data   map[string]interface{}
	Errors map[string][]string
}

// NewValidator returns a new Validator instance
func NewValidator() *Validator {
	return &Validator{
		Rules:  make(map[string]string),
		Data:   make(map[string]interface{}),
		Errors: make(map[string][]string),
	}
}

// AddRule adds a validation rule for a field
func (v *Validator) AddRule(field string, rule string) {
	v.Rules[field] = rule
}

// Validate validates the incoming data against the rules
func (v *Validator) Validate(c *fiber.Ctx) bool {
	var data map[string]interface{}
	//handle query parmas
	if c.Method() == "GET" || c.Method() == "DELETE" {
		queryData := c.Queries()
		data = make(map[string]interface{})
		for key, value := range queryData {
			data[key] = value
		}

	} else {

		// Handle multipart/form-data
		if strings.HasPrefix(c.Get("Content-Type"), "multipart/form-data") {
			form, err := c.MultipartForm()
			if err != nil {
				v.Errors["form"] = []string{"Invalid form data"}
				return false
			}
			data = make(map[string]interface{})
			for key, value := range form.Value {
				data[key] = value[0]
			}
			// Handle file fields separately
			for key, file := range form.File {
				data[key] = file[0].Filename
			}
		} else {
			//Handle application/json
			err := json.Unmarshal(c.Body(), &data)
			if err != nil {
				v.Errors["json"] = []string{"Invalid JSON"}
				return false
			}
		}
	}

	v.Data = data

	for field, rule := range v.Rules {
		value, ok := v.Data[field]

		// Handle the "present" rule first
		if strings.Contains(rule, "present") && !ok {
			v.Errors[field] = append(v.Errors[field], "This field must be present")
			continue
		}
		if !ok {
			if strings.Contains(rule, "nullable") {
				continue
			}
			v.Errors[field] = append(v.Errors[field], "This field is required")
			continue
		}

		rules := strings.Split(rule, "|")
		for _, r := range rules {
			switch {
			case strings.HasPrefix(r, "required"):
				if value == "" {
					v.Errors[field] = append(v.Errors[field], "This field is required")
				}
			case strings.HasPrefix(r, "string"):
				if value != nil {
					if reflect.TypeOf(value).Kind() != reflect.String {
						v.Errors[field] = append(v.Errors[field], "Must be a string")
					}
				}

			case strings.HasPrefix(r, "mobile"):
				pattern := `^[6-9]\d{9}$`
				if !regexp.MustCompile(pattern).MatchString(value.(string)) {
					v.Errors[field] = append(v.Errors[field], "Invalid format provided")
				}
			case strings.HasPrefix(r, "in:"):
				allowed := strings.Split(strings.TrimPrefix(r, "in:"), ",")
				if !contains(allowed, value.(string)) {
					v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Must be one of (%s)", strings.Join(allowed, ",")))
				}

			case strings.HasPrefix(r, "file"):
				fileHeader, err := c.FormFile(field)
				if err != nil || fileHeader == nil {
					// Ensure the field is either missing or empty
					if strings.Contains(rule, "nullable") {
						// If the field is missing or empty, skip validation
						if _, exists := v.Data[field]; !exists || v.Data[field] == "" {
							continue
						}
					}
					v.Errors[field] = append(v.Errors[field], "Must be a file")

				}
			case strings.HasPrefix(r, "docsize:"): //size in KB

				maxSizeKB := parseInt(strings.TrimPrefix(r, "docsize:"))
				maxSizeBytes := maxSizeKB * 1024 // Convert KB to bytes
				fileHeader, err := c.FormFile(field)
				if err != nil || fileHeader == nil {
					// Check if nullable is part of the rule
					if strings.Contains(rule, "nullable") {
						// If the field is missing or empty, skip validation
						if _, exists := v.Data[field]; !exists || v.Data[field] == "" {
							continue
						}
					}
					v.Errors[field] = append(v.Errors[field], "Error accessing file or no file")

				} else if fileHeader.Size > int64(maxSizeBytes) {
					v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Max file size %d kb", maxSizeKB))
				}
			case strings.HasPrefix(r, "email"):
				if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(value.(string)) {
					v.Errors[field] = append(v.Errors[field], "Invalid email")
				}
			case strings.HasPrefix(r, "min:"):
				min := strings.TrimPrefix(r, "min:")
				minValue, err := strconv.Atoi(min) // Convert the min value to an integer
				if err != nil {
					v.Errors[field] = append(v.Errors[field], "Invalid min value")
					break
				}

				switch val := value.(type) {
				case string:
					if len(val) < minValue {
						v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Min length %d", minValue))
					}
				case int:
					if val < minValue {
						v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Minimum value %d", minValue))
					}
				case float64:
					if int(val) < minValue {
						v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Minimum value %d", minValue))
					}
				default:
					v.Errors[field] = append(v.Errors[field], "Unsupported type for min validation")
				}

			case strings.HasPrefix(r, "max:"):
				max := strings.TrimPrefix(r, "max:")
				maxValue, err := strconv.Atoi(max) // Convert the max value to an integer
				if err != nil {
					v.Errors[field] = append(v.Errors[field], "Invalid max value")
					break
				}

				switch val := value.(type) {
				case string:
					if len(val) > maxValue {
						v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Max length %d", maxValue))
					}
				case int:
					if val > maxValue {
						v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Maximum value %d", maxValue))
					}
				case float64:
					if int(val) > maxValue {
						v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Maximum value %d", maxValue))
					}
				default:
					v.Errors[field] = append(v.Errors[field], "Unsupported type for max validation")
				}

			case strings.HasPrefix(r, "required_if:"):
				parts := strings.Split(strings.TrimPrefix(r, "required_if:"), ",")
				if len(parts) < 2 {
					v.Errors[field] = append(v.Errors[field], "Invalid required_if rule")
					break
				}
				otherField := parts[0]
				otherValue := parts[1]

				// Check if the other field exists
				otherFieldValue, exists := v.Data[otherField]
				if exists && fmt.Sprintf("%v", otherFieldValue) == otherValue {
					// Now check if the current field is present and not empty
					if value == nil || value == "" {
						v.Errors[field] = append(v.Errors[field], fmt.Sprintf("Required if %s is %s", otherField, otherValue))
					}
				}
			case strings.HasPrefix(r, "boolean"):
				if value != true && value != false {
					v.Errors[field] = append(v.Errors[field], "Must be a boolean")
				}
			case strings.HasPrefix(r, "array"):
				if reflect.TypeOf(value).Kind() != reflect.Slice {
					v.Errors[field] = append(v.Errors[field], "Must be an array")
				}

			case strings.HasPrefix(r, "numeric"):
				if strings.Contains(rule, "nullable") && value == nil {
					// Skip further validation if nullable and the value is nil or empty
					continue
				}
				if _, ok := value.(float64); !ok {
					v.Errors[field] = append(v.Errors[field], "Must be a number")
				}
			case strings.HasPrefix(r, "date"):
				if _, err := time.Parse("02-01-2006", value.(string)); err != nil { //day-month-yrr
					v.Errors[field] = append(v.Errors[field], "Invalid date format")
				}
			case strings.HasPrefix(r, "json"):

				_, ok := value.(map[string]interface{})
				if !ok {
					v.Errors[field] = append(v.Errors[field], "Must be a JSON object")
				}
			}
		}
	}

	return len(v.Errors) == 0
}

// GetErrors returns the validation errors
func (v *Validator) GetErrors() map[string][]string {
	return v.Errors
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
