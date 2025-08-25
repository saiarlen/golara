package database

import (
	"gorm.io/gorm"
)

// Model provides Laravel-style Eloquent model functionality
type Model struct {
	DB *gorm.DB
}

// NewModel creates a new model instance
func NewModel(db *gorm.DB) *Model {
	return &Model{DB: db}
}

// Query returns a new query builder
func (m *Model) Query() *QueryBuilder {
	return NewQueryBuilder(m.DB)
}

// Find finds a record by ID
func (m *Model) Find(id interface{}, dest interface{}) error {
	return m.DB.First(dest, id).Error
}

// FindOrFail finds a record by ID or returns error
func (m *Model) FindOrFail(id interface{}, dest interface{}) error {
	err := m.DB.First(dest, id).Error
	if err == gorm.ErrRecordNotFound {
		return &ModelNotFoundError{Model: dest, ID: id}
	}
	return err
}

// All returns all records
func (m *Model) All(dest interface{}) error {
	return m.DB.Find(dest).Error
}

// Create creates a new record
func (m *Model) Create(data interface{}) error {
	return m.DB.Create(data).Error
}

// Save saves a record
func (m *Model) Save(data interface{}) error {
	return m.DB.Save(data).Error
}

// Update updates a record
func (m *Model) Update(data interface{}) error {
	return m.DB.Updates(data).Error
}

// Delete deletes a record
func (m *Model) Delete(data interface{}) error {
	return m.DB.Delete(data).Error
}

// Where returns query builder with where condition
func (m *Model) Where(field string, operator string, value interface{}) *QueryBuilder {
	return m.Query().Where(field, operator, value)
}

// ModelNotFoundError represents a model not found error
type ModelNotFoundError struct {
	Model interface{}
	ID    interface{}
}

func (e *ModelNotFoundError) Error() string {
	return "record not found"
}