package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// QueryBuilder provides Laravel-style query building
type QueryBuilder struct {
	db       *gorm.DB
	model    interface{}
	query    *gorm.DB
	table    string
	selects  []string
	wheres   []string
	joins    []string
	orders   []string
	groups   []string
	havings  []string
	limit    int
	offset   int
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{
		db:    db,
		query: db,
	}
}

// Table sets the table name
func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.table = table
	qb.query = qb.db.Table(table)
	return qb
}

// Model sets the model for the query
func (qb *QueryBuilder) Model(model interface{}) *QueryBuilder {
	qb.model = model
	qb.query = qb.db.Model(model)
	return qb
}

// Select adds select fields
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.selects = append(qb.selects, fields...)
	qb.query = qb.query.Select(strings.Join(fields, ", "))
	return qb
}

// Where adds where condition
func (qb *QueryBuilder) Where(field string, operator string, value interface{}) *QueryBuilder {
	condition := fmt.Sprintf("%s %s ?", field, operator)
	qb.wheres = append(qb.wheres, condition)
	qb.query = qb.query.Where(condition, value)
	return qb
}

// WhereIn adds where in condition
func (qb *QueryBuilder) WhereIn(field string, values []interface{}) *QueryBuilder {
	qb.query = qb.query.Where(fmt.Sprintf("%s IN ?", field), values)
	return qb
}

// WhereNull adds where null condition
func (qb *QueryBuilder) WhereNull(field string) *QueryBuilder {
	qb.query = qb.query.Where(fmt.Sprintf("%s IS NULL", field))
	return qb
}

// WhereNotNull adds where not null condition
func (qb *QueryBuilder) WhereNotNull(field string) *QueryBuilder {
	qb.query = qb.query.Where(fmt.Sprintf("%s IS NOT NULL", field))
	return qb
}

// OrWhere adds or where condition
func (qb *QueryBuilder) OrWhere(field string, operator string, value interface{}) *QueryBuilder {
	condition := fmt.Sprintf("%s %s ?", field, operator)
	qb.query = qb.query.Or(condition, value)
	return qb
}

// Join adds join clause
func (qb *QueryBuilder) Join(table string, first string, operator string, second string) *QueryBuilder {
	joinClause := fmt.Sprintf("JOIN %s ON %s %s %s", table, first, operator, second)
	qb.joins = append(qb.joins, joinClause)
	qb.query = qb.query.Joins(fmt.Sprintf("%s ON %s %s %s", table, first, operator, second))
	return qb
}

// LeftJoin adds left join clause
func (qb *QueryBuilder) LeftJoin(table string, first string, operator string, second string) *QueryBuilder {
	qb.query = qb.query.Joins(fmt.Sprintf("LEFT JOIN %s ON %s %s %s", table, first, operator, second))
	return qb
}

// OrderBy adds order by clause
func (qb *QueryBuilder) OrderBy(field string, direction string) *QueryBuilder {
	order := fmt.Sprintf("%s %s", field, direction)
	qb.orders = append(qb.orders, order)
	qb.query = qb.query.Order(order)
	return qb
}

// GroupBy adds group by clause
func (qb *QueryBuilder) GroupBy(fields ...string) *QueryBuilder {
	qb.groups = append(qb.groups, fields...)
	qb.query = qb.query.Group(strings.Join(fields, ", "))
	return qb
}

// Having adds having clause
func (qb *QueryBuilder) Having(condition string, value interface{}) *QueryBuilder {
	qb.havings = append(qb.havings, condition)
	qb.query = qb.query.Having(condition, value)
	return qb
}

// Limit adds limit clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	qb.query = qb.query.Limit(limit)
	return qb
}

// Offset adds offset clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	qb.query = qb.query.Offset(offset)
	return qb
}

// Get executes the query and returns results
func (qb *QueryBuilder) Get(dest interface{}) error {
	return qb.query.Find(dest).Error
}

// First gets the first record
func (qb *QueryBuilder) First(dest interface{}) error {
	return qb.query.First(dest).Error
}

// Count returns the count of records
func (qb *QueryBuilder) Count() (int64, error) {
	var count int64
	err := qb.query.Count(&count).Error
	return count, err
}

// Paginate returns paginated results
func (qb *QueryBuilder) Paginate(page, perPage int, dest interface{}) (*PaginationResult, error) {
	offset := (page - 1) * perPage
	
	// Get total count
	total, err := qb.Count()
	if err != nil {
		return nil, err
	}
	
	// Get paginated data
	err = qb.Offset(offset).Limit(perPage).Get(dest)
	if err != nil {
		return nil, err
	}
	
	return &PaginationResult{
		Data:        dest,
		Total:       total,
		PerPage:     int64(perPage),
		CurrentPage: int64(page),
		LastPage:    (total + int64(perPage) - 1) / int64(perPage),
	}, nil
}

// Create inserts a new record
func (qb *QueryBuilder) Create(data interface{}) error {
	return qb.db.Create(data).Error
}

// Update updates records
func (qb *QueryBuilder) Update(data interface{}) error {
	return qb.query.Updates(data).Error
}

// Delete deletes records
func (qb *QueryBuilder) Delete() error {
	return qb.query.Delete(qb.model).Error
}

// PaginationResult represents paginated query results
type PaginationResult struct {
	Data        interface{} `json:"data"`
	Total       int64       `json:"total"`
	PerPage     int64       `json:"per_page"`
	CurrentPage int64       `json:"current_page"`
	LastPage    int64       `json:"last_page"`
	From        int64       `json:"from"`
	To          int64       `json:"to"`
}