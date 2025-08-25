package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestSuite provides testing utilities for Golara framework
type TestSuite struct {
	app *fiber.App
	t   *testing.T
}

// NewTestSuite creates a new test suite instance
func NewTestSuite(app *fiber.App, t *testing.T) *TestSuite {
	return &TestSuite{app: app, t: t}
}

// Get performs GET request for testing
func (ts *TestSuite) Get(url string) *TestResponse {
	req := httptest.NewRequest("GET", url, nil)
	resp, _ := ts.app.Test(req)
	return &TestResponse{Response: resp, t: ts.t}
}

// Post performs POST request with JSON body
func (ts *TestSuite) Post(url string, body interface{}) *TestResponse {
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := ts.app.Test(req)
	return &TestResponse{Response: resp, t: ts.t}
}

// TestResponse wraps HTTP response for testing assertions
type TestResponse struct {
	Response *http.Response
	t        *testing.T
}

// AssertStatus checks response status code
func (tr *TestResponse) AssertStatus(expectedStatus int) *TestResponse {
	if tr.Response.StatusCode != expectedStatus {
		tr.t.Errorf("Expected status %d, got %d", expectedStatus, tr.Response.StatusCode)
	}
	return tr
}

// AssertJSON checks if response contains expected JSON
func (tr *TestResponse) AssertJSON(expected map[string]interface{}) *TestResponse {
	body, _ := io.ReadAll(tr.Response.Body)
	var actual map[string]interface{}
	json.Unmarshal(body, &actual)
	
	for key, value := range expected {
		if actual[key] != value {
				tr.t.Errorf("Expected %s to be %v, got %v", key, value, actual[key])
		}
	}
	return tr
}