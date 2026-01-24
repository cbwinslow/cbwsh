package data

import (
	"testing"
)

func TestValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected DataType
	}{
		{"string value", "hello", DataTypeString},
		{"int value", 42, DataTypeInt},
		{"float value", 3.14, DataTypeFloat},
		{"bool value", true, DataTypeBool},
		{"null value", nil, DataTypeNull},
		{"array value", []string{"a", "b"}, DataTypeArray},
		{"object value", map[string]interface{}{"key": "value"}, DataTypeObject},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValue(tt.input)
			if v.Type != tt.expected {
				t.Errorf("Expected type %v, got %v", tt.expected, v.Type)
			}
		})
	}
}

func TestRecord(t *testing.T) {
	r := NewRecord()
	
	// Test Set and Get
	r.Set("name", "Alice")
	r.Set("age", 30)
	
	if !r.Has("name") {
		t.Error("Expected record to have 'name' field")
	}
	
	val, ok := r.Get("name")
	if !ok {
		t.Error("Expected to get 'name' field")
	}
	
	if val.String() != "Alice" {
		t.Errorf("Expected 'Alice', got %v", val.String())
	}
	
	// Test Has for non-existent field
	if r.Has("nonexistent") {
		t.Error("Expected record to not have 'nonexistent' field")
	}
}

func TestTable(t *testing.T) {
	table := NewTable()
	
	// Add records
	r1 := NewRecord()
	r1.Set("name", "Alice")
	r1.Set("age", 30)
	table.AddRecord(r1)
	
	r2 := NewRecord()
	r2.Set("name", "Bob")
	r2.Set("age", 25)
	table.AddRecord(r2)
	
	// Test length
	if table.Len() != 2 {
		t.Errorf("Expected 2 records, got %d", table.Len())
	}
	
	// Test columns
	if len(table.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(table.Columns))
	}
}

func TestTableFilter(t *testing.T) {
	table := NewTable()
	
	r1 := NewRecord()
	r1.Set("name", "Alice")
	r1.Set("age", 30)
	table.AddRecord(r1)
	
	r2 := NewRecord()
	r2.Set("name", "Bob")
	r2.Set("age", 25)
	table.AddRecord(r2)
	
	r3 := NewRecord()
	r3.Set("name", "Charlie")
	r3.Set("age", 35)
	table.AddRecord(r3)
	
	// Filter age > 26
	filtered := table.Filter(func(r *Record) bool {
		val, ok := r.Get("age")
		if !ok {
			return false
		}
		if age, ok := val.Data.(int); ok {
			return age > 26
		}
		return false
	})
	
	if filtered.Len() != 2 {
		t.Errorf("Expected 2 filtered records, got %d", filtered.Len())
	}
}

func TestTableSelect(t *testing.T) {
	table := NewTable()
	
	r1 := NewRecord()
	r1.Set("name", "Alice")
	r1.Set("age", 30)
	r1.Set("city", "NYC")
	table.AddRecord(r1)
	
	// Select only name and city
	selected := table.Select("name", "city")
	
	if len(selected.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(selected.Columns))
	}
	
	if selected.Len() != 1 {
		t.Errorf("Expected 1 record, got %d", selected.Len())
	}
	
	// Check that age is not in the selected table
	if selected.Records[0].Has("age") {
		t.Error("Expected 'age' to not be in selected records")
	}
}

func TestPipeline(t *testing.T) {
	table := NewTable()
	
	r1 := NewRecord()
	r1.Set("name", "Alice")
	r1.Set("age", 30)
	r1.Set("status", "active")
	table.AddRecord(r1)
	
	r2 := NewRecord()
	r2.Set("name", "Bob")
	r2.Set("age", 25)
	r2.Set("status", "inactive")
	table.AddRecord(r2)
	
	r3 := NewRecord()
	r3.Set("name", "Charlie")
	r3.Set("age", 35)
	r3.Set("status", "active")
	table.AddRecord(r3)
	
	// Pipeline: filter active users, select name and age, limit 1
	result := NewPipeline(table).
		Where("status", func(v *Value) bool {
			return v.String() == "active"
		}).
		Select("name", "age").
		Limit(1).
		Execute()
	
	if result.Len() != 1 {
		t.Errorf("Expected 1 record, got %d", result.Len())
	}
	
	if len(result.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(result.Columns))
	}
}
