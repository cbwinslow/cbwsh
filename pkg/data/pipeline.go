// Package data provides data pipeline operations
package data

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Pipeline represents a data transformation pipeline
type Pipeline struct {
	table      *Table
	operations []Operation
}

// Operation is a function that transforms a table
type Operation func(*Table) *Table

// NewPipeline creates a new data pipeline with the given table
func NewPipeline(table *Table) *Pipeline {
	return &Pipeline{
		table:      table,
		operations: []Operation{},
	}
}

// Where filters records based on a field condition
func (p *Pipeline) Where(field string, condition func(*Value) bool) *Pipeline {
	p.operations = append(p.operations, func(t *Table) *Table {
		return t.Filter(func(r *Record) bool {
			val, ok := r.Get(field)
			if !ok {
				return false
			}
			return condition(val)
		})
	})
	return p
}

// Select selects specific columns
func (p *Pipeline) Select(columns ...string) *Pipeline {
	p.operations = append(p.operations, func(t *Table) *Table {
		return t.Select(columns...)
	})
	return p
}

// Sort sorts records by a field
func (p *Pipeline) Sort(field string, ascending bool) *Pipeline {
	p.operations = append(p.operations, func(t *Table) *Table {
		result := NewTable()
		result.Columns = append([]string{}, t.Columns...)
		result.Records = append([]*Record{}, t.Records...)
		
		sort.Slice(result.Records, func(i, j int) bool {
			valI, okI := result.Records[i].Get(field)
			valJ, okJ := result.Records[j].Get(field)
			
			if !okI || !okJ {
				return false
			}
			
			// Compare based on type
			switch valI.Type {
			case DataTypeString:
				si, _ := valI.Data.(string)
				sj, _ := valJ.Data.(string)
				if ascending {
					return si < sj
				}
				return si > sj
			case DataTypeInt:
				ii, _ := valI.Data.(int)
				ij, _ := valJ.Data.(int)
				if ascending {
					return ii < ij
				}
				return ii > ij
			case DataTypeFloat:
				fi, _ := valI.Data.(float64)
				fj, _ := valJ.Data.(float64)
				if ascending {
					return fi < fj
				}
				return fi > fj
			}
			
			return false
		})
		
		return result
	})
	return p
}

// Limit limits the number of records
func (p *Pipeline) Limit(n int) *Pipeline {
	p.operations = append(p.operations, func(t *Table) *Table {
		result := NewTable()
		result.Columns = append([]string{}, t.Columns...)
		
		limit := n
		if limit > len(t.Records) {
			limit = len(t.Records)
		}
		
		result.Records = append([]*Record{}, t.Records[:limit]...)
		return result
	})
	return p
}

// GroupBy groups records by a field
func (p *Pipeline) GroupBy(field string) *Pipeline {
	p.operations = append(p.operations, func(t *Table) *Table {
		groups := make(map[string][]*Record)
		
		for _, record := range t.Records {
			val, ok := record.Get(field)
			if !ok {
				continue
			}
			
			key := val.String()
			groups[key] = append(groups[key], record)
		}
		
		result := NewTable()
		result.AddColumn(field)
		result.AddColumn("count")
		
		for key, records := range groups {
			newRecord := NewRecord()
			newRecord.Set(field, key)
			newRecord.Set("count", len(records))
			result.Records = append(result.Records, newRecord)
		}
		
		return result
	})
	return p
}

// Execute runs all operations in the pipeline
func (p *Pipeline) Execute() *Table {
	result := p.table
	for _, op := range p.operations {
		result = op(result)
	}
	return result
}

// Formatter formats a table for display
type Formatter interface {
	Format(*Table) string
}

// TableFormatter formats a table as a text table
type TableFormatter struct {
	MaxCellWidth int
}

// NewTableFormatter creates a new table formatter
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{
		MaxCellWidth: 50,
	}
}

// Format formats a table as a text table
func (f *TableFormatter) Format(table *Table) string {
	if len(table.Records) == 0 {
		return "No data"
	}
	
	// Calculate column widths
	widths := make(map[string]int)
	for _, col := range table.Columns {
		widths[col] = len(col)
	}
	
	for _, record := range table.Records {
		for _, col := range table.Columns {
			if val, ok := record.Get(col); ok {
				length := len(val.String())
				if length > widths[col] {
					widths[col] = length
				}
				if widths[col] > f.MaxCellWidth {
					widths[col] = f.MaxCellWidth
				}
			}
		}
	}
	
	var sb strings.Builder
	
	// Header
	sb.WriteString("┌")
	for i, col := range table.Columns {
		sb.WriteString(strings.Repeat("─", widths[col]+2))
		if i < len(table.Columns)-1 {
			sb.WriteString("┬")
		}
	}
	sb.WriteString("┐\n")
	
	// Column names
	sb.WriteString("│")
	for _, col := range table.Columns {
		sb.WriteString(" ")
		sb.WriteString(padRight(col, widths[col]))
		sb.WriteString(" │")
	}
	sb.WriteString("\n")
	
	// Separator
	sb.WriteString("├")
	for i, col := range table.Columns {
		sb.WriteString(strings.Repeat("─", widths[col]+2))
		if i < len(table.Columns)-1 {
			sb.WriteString("┼")
		}
	}
	sb.WriteString("┤\n")
	
	// Data rows
	for _, record := range table.Records {
		sb.WriteString("│")
		for _, col := range table.Columns {
			sb.WriteString(" ")
			if val, ok := record.Get(col); ok {
				text := val.String()
				if len(text) > f.MaxCellWidth {
					text = text[:f.MaxCellWidth-3] + "..."
				}
				sb.WriteString(padRight(text, widths[col]))
			} else {
				sb.WriteString(padRight("", widths[col]))
			}
			sb.WriteString(" │")
		}
		sb.WriteString("\n")
	}
	
	// Footer
	sb.WriteString("└")
	for i, col := range table.Columns {
		sb.WriteString(strings.Repeat("─", widths[col]+2))
		if i < len(table.Columns)-1 {
			sb.WriteString("┴")
		}
	}
	sb.WriteString("┘\n")
	
	sb.WriteString(fmt.Sprintf("\n%d rows\n", len(table.Records)))
	
	return sb.String()
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

// JSONFormatter formats a table as JSON
type JSONFormatter struct {
	Pretty bool
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter(pretty bool) *JSONFormatter {
	return &JSONFormatter{Pretty: pretty}
}

// Format formats a table as JSON
func (f *JSONFormatter) Format(table *Table) string {
	var result []map[string]interface{}
	
	for _, record := range table.Records {
		result = append(result, record.ToMap())
	}
	
	var data []byte
	var err error
	
	if f.Pretty {
		data, err = jsonMarshalIndent(result, "", "  ")
	} else {
		data, err = jsonMarshal(result)
	}
	
	if err != nil {
		return "{}"
	}
	
	return string(data)
}

// Helper functions for JSON formatting
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func jsonMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
