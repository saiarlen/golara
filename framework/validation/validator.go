package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Validator provides enhanced validation functionality
type Validator struct {
	rules      map[string][]Rule
	data       map[string]interface{}
	errors     map[string][]string
	customRules map[string]CustomRule
}

// Rule represents a validation rule
type Rule interface {
	Validate(field string, value interface{}, data map[string]interface{}) error
	GetMessage() string
}

// CustomRule represents a custom validation rule
type CustomRule func(field string, value interface{}, parameters []string, data map[string]interface{}) error

// NewValidator creates a new validator
func NewValidator() *Validator {
	v := &Validator{
		rules:       make(map[string][]Rule),
		data:        make(map[string]interface{}),
		errors:      make(map[string][]string),
		customRules: make(map[string]CustomRule),
	}
	
	// Register default custom rules
	v.registerDefaultRules()
	
	return v
}

// AddRule adds a validation rule for a field
func (v *Validator) AddRule(field string, rule Rule) *Validator {
	v.rules[field] = append(v.rules[field], rule)
	return v
}

// AddCustomRule registers a custom validation rule
func (v *Validator) AddCustomRule(name string, rule CustomRule) {
	v.customRules[name] = rule
}

// Validate validates data against rules
func (v *Validator) Validate(c *fiber.Ctx) bool {
	v.parseRequestData(c)
	
	for field, rules := range v.rules {
		value, exists := v.data[field]
		
		for _, rule := range rules {
			if err := rule.Validate(field, value, v.data); err != nil {
				v.addError(field, err.Error())
			}
		}
		
		// Check if field is required but missing
		if !exists {
			for _, rule := range rules {
				if _, ok := rule.(*RequiredRule); ok {
					v.addError(field, "This field is required")
					break
				}
			}
		}
	}
	
	return len(v.errors) == 0
}

// GetErrors returns validation errors
func (v *Validator) GetErrors() map[string][]string {
	return v.errors
}

// GetData returns parsed data
func (v *Validator) GetData() map[string]interface{} {
	return v.data
}

func (v *Validator) addError(field, message string) {
	v.errors[field] = append(v.errors[field], message)
}

func (v *Validator) parseRequestData(c *fiber.Ctx) {
	var data map[string]interface{}
	
	if c.Method() == "GET" || c.Method() == "DELETE" {
		// Parse query parameters
		data = make(map[string]interface{})
		for key, value := range c.Queries() {
			data[key] = value
		}
	} else {
		// Parse JSON body
		if err := c.BodyParser(&data); err != nil {
			// Try form data
			form, err := c.MultipartForm()
			if err == nil {
				data = make(map[string]interface{})
				for key, values := range form.Value {
					if len(values) > 0 {
						data[key] = values[0]
					}
				}
			}
		}
	}
	
	v.data = data
}

// Built-in validation rules

// RequiredRule validates required fields
type RequiredRule struct{}

func (r *RequiredRule) Validate(field string, value interface{}, data map[string]interface{}) error {
	if value == nil || value == "" {
		return fmt.Errorf("The %s field is required", field)
	}
	return nil
}

func (r *RequiredRule) GetMessage() string {
	return "This field is required"
}

// EmailRule validates email format
type EmailRule struct{}

func (r *EmailRule) Validate(field string, value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil // Skip if nil (use Required rule for required validation)
	}
	
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("The %s field must be a string", field)
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(str) {
		return fmt.Errorf("The %s field must be a valid email address", field)
	}
	
	return nil
}

func (r *EmailRule) GetMessage() string {
	return "Must be a valid email address"
}

// MinRule validates minimum length/value
type MinRule struct {
	Min int
}

func (r *MinRule) Validate(field string, value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case string:
		if len(v) < r.Min {
			return fmt.Errorf("The %s field must be at least %d characters", field, r.Min)
		}
	case int:
		if v < r.Min {
			return fmt.Errorf("The %s field must be at least %d", field, r.Min)
		}
	case float64:
		if int(v) < r.Min {
			return fmt.Errorf("The %s field must be at least %d", field, r.Min)
		}
	default:
		return fmt.Errorf("The %s field type is not supported for min validation", field)
	}
	
	return nil
}

func (r *MinRule) GetMessage() string {
	return fmt.Sprintf("Must be at least %d", r.Min)
}

// MaxRule validates maximum length/value
type MaxRule struct {
	Max int
}

func (r *MaxRule) Validate(field string, value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}
	
	switch v := value.(type) {
	case string:
		if len(v) > r.Max {
			return fmt.Errorf("The %s field must not exceed %d characters", field, r.Max)
		}
	case int:
		if v > r.Max {
			return fmt.Errorf("The %s field must not exceed %d", field, r.Max)
		}
	case float64:
		if int(v) > r.Max {
			return fmt.Errorf("The %s field must not exceed %d", field, r.Max)
		}
	default:
		return fmt.Errorf("The %s field type is not supported for max validation", field)
	}
	
	return nil
}

func (r *MaxRule) GetMessage() string {
	return fmt.Sprintf("Must not exceed %d", r.Max)
}

// InRule validates value is in allowed list
type InRule struct {
	Values []string
}

func (r *InRule) Validate(field string, value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}
	
	str := fmt.Sprintf("%v", value)
	for _, allowed := range r.Values {
		if str == allowed {
			return nil
		}
	}
	
	return fmt.Errorf("The %s field must be one of: %s", field, strings.Join(r.Values, ", "))
}

func (r *InRule) GetMessage() string {
	return "Must be one of the allowed values"
}

// NumericRule validates numeric values
type NumericRule struct{}

func (r *NumericRule) Validate(field string, value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}
	
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return nil
	case string:
		if _, err := strconv.ParseFloat(value.(string), 64); err != nil {
			return fmt.Errorf("The %s field must be numeric", field)
		}
		return nil
	default:
		return fmt.Errorf("The %s field must be numeric", field)
	}
}

func (r *NumericRule) GetMessage() string {
	return "Must be numeric"
}

// DateRule validates date format
type DateRule struct {
	Format string
}

func (r *DateRule) Validate(field string, value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}
	
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("The %s field must be a string", field)
	}
	
	format := r.Format
	if format == "" {
		format = "2006-01-02" // Default format
	}
	
	if _, err := time.Parse(format, str); err != nil {
		return fmt.Errorf("The %s field must be a valid date in format %s", field, format)
	}
	
	return nil
}

func (r *DateRule) GetMessage() string {
	return "Must be a valid date"
}

// Helper functions for creating rules
func Required() Rule {
	return &RequiredRule{}
}

func Email() Rule {
	return &EmailRule{}
}

func Min(min int) Rule {
	return &MinRule{Min: min}
}

func Max(max int) Rule {
	return &MaxRule{Max: max}
}

func In(values ...string) Rule {
	return &InRule{Values: values}
}

func Numeric() Rule {
	return &NumericRule{}
}

func Date(format ...string) Rule {
	f := "2006-01-02"
	if len(format) > 0 {
		f = format[0]
	}
	return &DateRule{Format: f}
}

func (v *Validator) registerDefaultRules() {
	// Register some common custom rules
	v.AddCustomRule("confirmed", func(field string, value interface{}, parameters []string, data map[string]interface{}) error {
		confirmField := field + "_confirmation"
		if len(parameters) > 0 {
			confirmField = parameters[0]
		}
		
		confirmValue, exists := data[confirmField]
		if !exists {
			return fmt.Errorf("The %s confirmation field is missing", field)
		}
		
		if value != confirmValue {
			return fmt.Errorf("The %s confirmation does not match", field)
		}
		
		return nil
	})
	
	v.AddCustomRule("unique", func(field string, value interface{}, parameters []string, data map[string]interface{}) error {
		// This would typically check database uniqueness
		// Implementation depends on your database setup
		return nil
	})
}