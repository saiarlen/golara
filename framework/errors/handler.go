package errors

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	ValidationError ErrorType = "validation"
	DatabaseError   ErrorType = "database"
	AuthError       ErrorType = "authentication"
	NotFoundError   ErrorType = "not_found"
	InternalError   ErrorType = "internal"
)

// AppError represents a structured application error
type AppError struct {
	Type       ErrorType              `json:"type"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
	Timestamp  time.Time              `json:"timestamp"`
	RequestID  string                 `json:"request_id,omitempty"`
	Stack      string                 `json:"stack,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(errType ErrorType, message string, statusCode int) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		StatusCode: statusCode,
		Timestamp:  time.Now(),
		Details:    make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithStack adds stack trace to the error
func (e *AppError) WithStack() *AppError {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	e.Stack = string(buf[:n])
	return e
}

// ErrorHandler is the global error handler for the framework
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Set request ID for tracking
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		requestID = generateRequestID()
		c.Set("X-Request-ID", requestID)
	}

	var appErr *AppError
	
	// Check if it's already an AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
		appErr.RequestID = requestID
	} else if e, ok := err.(*fiber.Error); ok {
		// Handle Fiber errors
		appErr = &AppError{
			Type:       InternalError,
			Message:    e.Message,
			StatusCode: e.Code,
			Timestamp:  time.Now(),
			RequestID:  requestID,
		}
	} else {
		// Handle generic errors
		appErr = &AppError{
			Type:       InternalError,
			Message:    "Internal server error",
			StatusCode: 500,
			Timestamp:  time.Now(),
			RequestID:  requestID,
			Details:    map[string]interface{}{"original_error": err.Error()},
		}
	}

	// Log the error
	logError(appErr, c)

	// Return appropriate response based on environment
	if isProduction() {
		return c.Status(appErr.StatusCode).JSON(sanitizeError(appErr))
	}
	
	return c.Status(appErr.StatusCode).JSON(appErr)
}

// Helper functions for common errors
func ValidationErr(message string, details map[string]interface{}) *AppError {
	return NewAppError(ValidationError, message, 400).WithDetails(details)
}

func NotFoundErr(message string) *AppError {
	return NewAppError(NotFoundError, message, 404)
}

func AuthErr(message string) *AppError {
	return NewAppError(AuthError, message, 401)
}

func InternalErr(message string) *AppError {
	return NewAppError(InternalError, message, 500).WithStack()
}

func DatabaseErr(message string, details map[string]interface{}) *AppError {
	return NewAppError(DatabaseError, message, 500).WithDetails(details)
}

func logError(err *AppError, c *fiber.Ctx) {
	log.Printf("[ERROR] %s | %s %s | %s | %s | %v",
		err.RequestID,
		c.Method(),
		c.Path(),
		c.IP(),
		err.Type,
		err.Message,
	)
	
	if err.Stack != "" {
		log.Printf("[STACK] %s | %s", err.RequestID, err.Stack)
	}
}

func sanitizeError(err *AppError) map[string]interface{} {
	return map[string]interface{}{
		"type":       err.Type,
		"message":    err.Message,
		"timestamp":  err.Timestamp,
		"request_id": err.RequestID,
	}
}

func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func isProduction() bool {
	// This should check your environment configuration
	return false // For now, always return false for development
}