// Code generated by ent, DO NOT EDIT.

package process

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the process type in the database.
	Label = "process"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldFilename holds the string denoting the filename field in the database.
	FieldFilename = "filename"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldHash holds the string denoting the hash field in the database.
	FieldHash = "hash"
	// FieldProcessResponse holds the string denoting the process_response field in the database.
	FieldProcessResponse = "process_response"
	// FieldAnchorID holds the string denoting the anchor_id field in the database.
	FieldAnchorID = "anchor_id"
	// FieldIsAggregated holds the string denoting the is_aggregated field in the database.
	FieldIsAggregated = "is_aggregated"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// Table holds the table name of the process in the database.
	Table = "processes"
)

// Columns holds all SQL columns for process fields.
var Columns = []string{
	FieldID,
	FieldFilename,
	FieldStatus,
	FieldHash,
	FieldProcessResponse,
	FieldAnchorID,
	FieldIsAggregated,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// FilenameValidator is a validator for the "filename" field. It is called by the builders before save.
	FilenameValidator func(string) error
	// DefaultStatus holds the default value on creation for the "status" field.
	DefaultStatus bool
	// HashValidator is a validator for the "hash" field. It is called by the builders before save.
	HashValidator func(string) error
	// DefaultIsAggregated holds the default value on creation for the "is_aggregated" field.
	DefaultIsAggregated bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Process queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByFilename orders the results by the filename field.
func ByFilename(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFilename, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByHash orders the results by the hash field.
func ByHash(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHash, opts...).ToFunc()
}

// ByAnchorID orders the results by the anchor_id field.
func ByAnchorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAnchorID, opts...).ToFunc()
}

// ByIsAggregated orders the results by the is_aggregated field.
func ByIsAggregated(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsAggregated, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}