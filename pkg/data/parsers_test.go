package data

import (
	"strings"
	"testing"
)

func TestJSONParser(t *testing.T) {
	parser := NewJSONParser()
	
	jsonData := `[
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25}
	]`
	
	table, err := parser.Parse([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	if table.Len() != 2 {
		t.Errorf("Expected 2 records, got %d", table.Len())
	}
	
	if len(table.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(table.Columns))
	}
}

func TestJSONParserSingleObject(t *testing.T) {
	parser := NewJSONParser()
	
	jsonData := `{"name": "Alice", "age": 30}`
	
	table, err := parser.Parse([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	if table.Len() != 1 {
		t.Errorf("Expected 1 record, got %d", table.Len())
	}
}

func TestYAMLParser(t *testing.T) {
	parser := NewYAMLParser()
	
	yamlData := `
- name: Alice
  age: 30
- name: Bob
  age: 25
`
	
	table, err := parser.Parse([]byte(yamlData))
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}
	
	if table.Len() != 2 {
		t.Errorf("Expected 2 records, got %d", table.Len())
	}
}

func TestCSVParser(t *testing.T) {
	parser := NewCSVParser()
	
	csvData := `name,age
Alice,30
Bob,25
`
	
	table, err := parser.Parse([]byte(csvData))
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}
	
	if table.Len() != 2 {
		t.Errorf("Expected 2 records, got %d", table.Len())
	}
	
	if len(table.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(table.Columns))
	}
}

func TestCSVParserNoHeader(t *testing.T) {
	parser := NewCSVParser()
	parser.HasHeader = false
	
	csvData := `Alice,30
Bob,25
`
	
	table, err := parser.Parse([]byte(csvData))
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}
	
	if table.Len() != 2 {
		t.Errorf("Expected 2 records, got %d", table.Len())
	}
	
	// Should have generated column names
	if len(table.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(table.Columns))
	}
}

func TestGetParser(t *testing.T) {
	tests := []struct {
		format   string
		expected bool
	}{
		{"json", true},
		{"yaml", true},
		{"yml", true},
		{"csv", true},
		{"xml", false},
		{"", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			parser := GetParser(tt.format)
			hasParser := parser != nil
			if hasParser != tt.expected {
				t.Errorf("Expected parser existence %v for format %s, got %v", 
					tt.expected, tt.format, hasParser)
			}
		})
	}
}

func TestJSONParserWithReader(t *testing.T) {
	parser := NewJSONParser()
	
	jsonData := `[{"name": "Alice"}]`
	reader := strings.NewReader(jsonData)
	
	table, err := parser.ParseReader(reader)
	if err != nil {
		t.Fatalf("Failed to parse JSON from reader: %v", err)
	}
	
	if table.Len() != 1 {
		t.Errorf("Expected 1 record, got %d", table.Len())
	}
}

func TestCSVParserWithReader(t *testing.T) {
	parser := NewCSVParser()
	
	csvData := `name,age
Alice,30`
	reader := strings.NewReader(csvData)
	
	table, err := parser.ParseReader(reader)
	if err != nil {
		t.Fatalf("Failed to parse CSV from reader: %v", err)
	}
	
	if table.Len() != 1 {
		t.Errorf("Expected 1 record, got %d", table.Len())
	}
}
