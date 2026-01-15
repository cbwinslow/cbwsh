// Package data provides structured data processing capabilities for cbwsh.
// Inspired by modern shells like Nushell and PowerShell, this package enables
// type-safe data pipelines and native support for common data formats.
package data

import (
	"encoding/json"
	"time"
)

// DataType represents the type of structured data
type DataType string

const (
	// DataTypeString represents string data
	DataTypeString DataType = "string"
	// DataTypeInt represents integer data
	DataTypeInt DataType = "int"
	// DataTypeFloat represents floating point data
	DataTypeFloat DataType = "float"
	// DataTypeBool represents boolean data
	DataTypeBool DataType = "bool"
	// DataTypeArray represents array/list data
	DataTypeArray DataType = "array"
	// DataTypeObject represents object/map data
	DataTypeObject DataType = "object"
	// DataTypeNull represents null/nil data
	DataTypeNull DataType = "null"
	// DataTypeTime represents time/date data
	DataTypeTime DataType = "time"
)

// Value represents a typed value in the data pipeline
type Value struct {
	Type  DataType    `json:"type"`
	Data  interface{} `json:"data"`
	raw   []byte      // Original raw data for caching
}

// NewValue creates a new typed value
func NewValue(data interface{}) *Value {
	return &Value{
		Type: inferType(data),
		Data: data,
	}
}

// String returns the string representation of the value
func (v *Value) String() string {
	if v.Data == nil {
		return ""
	}
	
	switch v.Type {
	case DataTypeString:
		if s, ok := v.Data.(string); ok {
			return s
		}
	case DataTypeBool:
		if b, ok := v.Data.(bool); ok {
			if b {
				return "true"
			}
			return "false"
		}
	}
	
	// For complex types, use JSON encoding
	data, err := json.Marshal(v.Data)
	if err != nil {
		return ""
	}
	return string(data)
}

// Record represents a single record in a structured data stream
// Similar to a row in a table or a JSON object
type Record struct {
	Fields map[string]*Value `json:"fields"`
}

// NewRecord creates a new record
func NewRecord() *Record {
	return &Record{
		Fields: make(map[string]*Value),
	}
}

// Set sets a field value in the record
func (r *Record) Set(key string, value interface{}) {
	r.Fields[key] = NewValue(value)
}

// Get gets a field value from the record
func (r *Record) Get(key string) (*Value, bool) {
	val, ok := r.Fields[key]
	return val, ok
}

// Has checks if a field exists in the record
func (r *Record) Has(key string) bool {
	_, ok := r.Fields[key]
	return ok
}

// ToMap converts the record to a map
func (r *Record) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range r.Fields {
		result[k] = v.Data
	}
	return result
}

// Table represents a collection of records (like a database table)
type Table struct {
	Columns []string  `json:"columns"`
	Records []*Record `json:"records"`
}

// NewTable creates a new table
func NewTable() *Table {
	return &Table{
		Columns: []string{},
		Records: []*Record{},
	}
}

// AddColumn adds a column to the table
func (t *Table) AddColumn(name string) {
	// Check if column already exists
	for _, col := range t.Columns {
		if col == name {
			return
		}
	}
	t.Columns = append(t.Columns, name)
}

// AddRecord adds a record to the table
func (t *Table) AddRecord(record *Record) {
	// Ensure all fields are in columns
	for key := range record.Fields {
		t.AddColumn(key)
	}
	t.Records = append(t.Records, record)
}

// Filter filters records based on a predicate function
func (t *Table) Filter(predicate func(*Record) bool) *Table {
	result := NewTable()
	result.Columns = append([]string{}, t.Columns...)
	
	for _, record := range t.Records {
		if predicate(record) {
			result.Records = append(result.Records, record)
		}
	}
	
	return result
}

// Select selects specific columns from the table
func (t *Table) Select(columns ...string) *Table {
	result := NewTable()
	result.Columns = columns
	
	for _, record := range t.Records {
		newRecord := NewRecord()
		for _, col := range columns {
			if val, ok := record.Get(col); ok {
				newRecord.Fields[col] = val
			}
		}
		result.Records = append(result.Records, newRecord)
	}
	
	return result
}

// Len returns the number of records
func (t *Table) Len() int {
	return len(t.Records)
}

// inferType infers the DataType from a Go value
func inferType(data interface{}) DataType {
	if data == nil {
		return DataTypeNull
	}
	
	switch data.(type) {
	case string:
		return DataTypeString
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return DataTypeInt
	case float32, float64:
		return DataTypeFloat
	case bool:
		return DataTypeBool
	case []interface{}, []string, []int, []float64, []bool:
		return DataTypeArray
	case map[string]interface{}:
		return DataTypeObject
	case time.Time:
		return DataTypeTime
	default:
		return DataTypeObject
	}
}
