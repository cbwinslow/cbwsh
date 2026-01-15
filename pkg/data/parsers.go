// Package data provides parsers for common data formats
package data

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parser is an interface for parsing data formats
type Parser interface {
	Parse(data []byte) (*Table, error)
	ParseReader(r io.Reader) (*Table, error)
}

// JSONParser parses JSON data into structured tables
type JSONParser struct{}

// NewJSONParser creates a new JSON parser
func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

// Parse parses JSON bytes into a table
func (p *JSONParser) Parse(data []byte) (*Table, error) {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	
	return p.convertToTable(raw), nil
}

// ParseReader parses JSON from a reader into a table
func (p *JSONParser) ParseReader(r io.Reader) (*Table, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return p.Parse(data)
}

func (p *JSONParser) convertToTable(data interface{}) *Table {
	table := NewTable()
	
	switch v := data.(type) {
	case []interface{}:
		// Array of objects
		for _, item := range v {
			if obj, ok := item.(map[string]interface{}); ok {
				record := NewRecord()
				for key, val := range obj {
					record.Set(key, val)
				}
				table.AddRecord(record)
			}
		}
	case map[string]interface{}:
		// Single object - convert to single-row table
		record := NewRecord()
		for key, val := range v {
			record.Set(key, val)
		}
		table.AddRecord(record)
	}
	
	return table
}

// YAMLParser parses YAML data into structured tables
type YAMLParser struct{}

// NewYAMLParser creates a new YAML parser
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

// Parse parses YAML bytes into a table
func (p *YAMLParser) Parse(data []byte) (*Table, error) {
	var raw interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	
	return p.convertToTable(raw), nil
}

// ParseReader parses YAML from a reader into a table
func (p *YAMLParser) ParseReader(r io.Reader) (*Table, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return p.Parse(data)
}

func (p *YAMLParser) convertToTable(data interface{}) *Table {
	table := NewTable()
	
	switch v := data.(type) {
	case []interface{}:
		// Array of objects
		for _, item := range v {
			if obj, ok := item.(map[string]interface{}); ok {
				record := NewRecord()
				for key, val := range obj {
					record.Set(key, val)
				}
				table.AddRecord(record)
			}
		}
	case map[string]interface{}:
		// Single object - convert to single-row table
		record := NewRecord()
		for key, val := range v {
			record.Set(key, val)
		}
		table.AddRecord(record)
	}
	
	return table
}

// CSVParser parses CSV data into structured tables
type CSVParser struct {
	Delimiter rune
	HasHeader bool
}

// NewCSVParser creates a new CSV parser
func NewCSVParser() *CSVParser {
	return &CSVParser{
		Delimiter: ',',
		HasHeader: true,
	}
}

// Parse parses CSV bytes into a table
func (p *CSVParser) Parse(data []byte) (*Table, error) {
	return p.ParseReader(strings.NewReader(string(data)))
}

// ParseReader parses CSV from a reader into a table
func (p *CSVParser) ParseReader(r io.Reader) (*Table, error) {
	reader := csv.NewReader(r)
	reader.Comma = p.Delimiter
	
	table := NewTable()
	
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	
	if len(records) == 0 {
		return table, nil
	}
	
	var headers []string
	startIdx := 0
	
	if p.HasHeader && len(records) > 0 {
		headers = records[0]
		startIdx = 1
	} else {
		// Generate column names: col0, col1, col2, etc.
		for i := 0; i < len(records[0]); i++ {
			headers = append(headers, string(rune('A'+i)))
		}
	}
	
	table.Columns = headers
	
	// Parse data rows
	for _, row := range records[startIdx:] {
		record := NewRecord()
		for i, val := range row {
			if i < len(headers) {
				record.Set(headers[i], val)
			}
		}
		table.Records = append(table.Records, record)
	}
	
	return table, nil
}

// GetParser returns the appropriate parser for a given format
func GetParser(format string) Parser {
	switch strings.ToLower(format) {
	case "json":
		return NewJSONParser()
	case "yaml", "yml":
		return NewYAMLParser()
	case "csv":
		return NewCSVParser()
	default:
		return nil
	}
}
