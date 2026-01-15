// Package data provides command implementations for structured data operations
package data

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cbwinslow/cbwsh/pkg/data"
)

// ParseCommand parses structured data from stdin or files
type ParseCommand struct {
	Format string
	Input  io.Reader
}

// NewParseCommand creates a new parse command
func NewParseCommand(format string, input io.Reader) *ParseCommand {
	return &ParseCommand{
		Format: format,
		Input:  input,
	}
}

// Execute runs the parse command
func (c *ParseCommand) Execute() (*data.Table, error) {
	parser := data.GetParser(c.Format)
	if parser == nil {
		return nil, fmt.Errorf("unsupported format: %s", c.Format)
	}
	
	return parser.ParseReader(c.Input)
}

// WhereCommand filters table data
type WhereCommand struct {
	Table     *data.Table
	Field     string
	Operator  string
	Value     string
}

// NewWhereCommand creates a new where command
func NewWhereCommand(table *data.Table, field, operator, value string) *WhereCommand {
	return &WhereCommand{
		Table:    table,
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// Execute runs the where command
func (c *WhereCommand) Execute() (*data.Table, error) {
	pipeline := data.NewPipeline(c.Table)
	
	switch c.Operator {
	case "==", "=":
		pipeline = pipeline.Where(c.Field, func(v *data.Value) bool {
			return v.String() == c.Value
		})
	case "!=":
		pipeline = pipeline.Where(c.Field, func(v *data.Value) bool {
			return v.String() != c.Value
		})
	case "contains":
		pipeline = pipeline.Where(c.Field, func(v *data.Value) bool {
			return strings.Contains(v.String(), c.Value)
		})
	default:
		return nil, fmt.Errorf("unsupported operator: %s", c.Operator)
	}
	
	return pipeline.Execute(), nil
}

// SelectCommand selects specific columns
type SelectCommand struct {
	Table   *data.Table
	Columns []string
}

// NewSelectCommand creates a new select command
func NewSelectCommand(table *data.Table, columns []string) *SelectCommand {
	return &SelectCommand{
		Table:   table,
		Columns: columns,
	}
}

// Execute runs the select command
func (c *SelectCommand) Execute() (*data.Table, error) {
	pipeline := data.NewPipeline(c.Table)
	return pipeline.Select(c.Columns...).Execute(), nil
}

// SortCommand sorts table data
type SortCommand struct {
	Table      *data.Table
	Field      string
	Descending bool
}

// NewSortCommand creates a new sort command
func NewSortCommand(table *data.Table, field string, descending bool) *SortCommand {
	return &SortCommand{
		Table:      table,
		Field:      field,
		Descending: descending,
	}
}

// Execute runs the sort command
func (c *SortCommand) Execute() (*data.Table, error) {
	pipeline := data.NewPipeline(c.Table)
	return pipeline.Sort(c.Field, !c.Descending).Execute(), nil
}

// LimitCommand limits the number of rows
type LimitCommand struct {
	Table *data.Table
	Count int
}

// NewLimitCommand creates a new limit command
func NewLimitCommand(table *data.Table, count int) *LimitCommand {
	return &LimitCommand{
		Table: table,
		Count: count,
	}
}

// Execute runs the limit command
func (c *LimitCommand) Execute() (*data.Table, error) {
	pipeline := data.NewPipeline(c.Table)
	return pipeline.Limit(c.Count).Execute(), nil
}

// GroupByCommand groups data by a field
type GroupByCommand struct {
	Table *data.Table
	Field string
}

// NewGroupByCommand creates a new group-by command
func NewGroupByCommand(table *data.Table, field string) *GroupByCommand {
	return &GroupByCommand{
		Table: table,
		Field: field,
	}
}

// Execute runs the group-by command
func (c *GroupByCommand) Execute() (*data.Table, error) {
	pipeline := data.NewPipeline(c.Table)
	return pipeline.GroupBy(c.Field).Execute(), nil
}

// FormatCommand formats table data for output
type FormatCommand struct {
	Table      *data.Table
	OutputType string
}

// NewFormatCommand creates a new format command
func NewFormatCommand(table *data.Table, outputType string) *FormatCommand {
	return &FormatCommand{
		Table:      table,
		OutputType: outputType,
	}
}

// Execute runs the format command
func (c *FormatCommand) Execute() (string, error) {
	switch c.OutputType {
	case "table", "":
		formatter := data.NewTableFormatter()
		return formatter.Format(c.Table), nil
	case "json":
		formatter := data.NewJSONFormatter(true)
		return formatter.Format(c.Table), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", c.OutputType)
	}
}

// CommandRunner executes a chain of data commands
type CommandRunner struct {
	table *data.Table
}

// NewCommandRunner creates a new command runner
func NewCommandRunner() *CommandRunner {
	return &CommandRunner{}
}

// Parse parses data from a source
func (r *CommandRunner) Parse(format string, source io.Reader) error {
	cmd := NewParseCommand(format, source)
	table, err := cmd.Execute()
	if err != nil {
		return err
	}
	r.table = table
	return nil
}

// Where filters the data
func (r *CommandRunner) Where(field, operator, value string) error {
	if r.table == nil {
		return fmt.Errorf("no data loaded")
	}
	cmd := NewWhereCommand(r.table, field, operator, value)
	table, err := cmd.Execute()
	if err != nil {
		return err
	}
	r.table = table
	return nil
}

// Select selects columns
func (r *CommandRunner) Select(columns ...string) error {
	if r.table == nil {
		return fmt.Errorf("no data loaded")
	}
	cmd := NewSelectCommand(r.table, columns)
	table, err := cmd.Execute()
	if err != nil {
		return err
	}
	r.table = table
	return nil
}

// Sort sorts the data
func (r *CommandRunner) Sort(field string, descending bool) error {
	if r.table == nil {
		return fmt.Errorf("no data loaded")
	}
	cmd := NewSortCommand(r.table, field, descending)
	table, err := cmd.Execute()
	if err != nil {
		return err
	}
	r.table = table
	return nil
}

// Limit limits the number of rows
func (r *CommandRunner) Limit(count int) error {
	if r.table == nil {
		return fmt.Errorf("no data loaded")
	}
	cmd := NewLimitCommand(r.table, count)
	table, err := cmd.Execute()
	if err != nil {
		return err
	}
	r.table = table
	return nil
}

// GroupBy groups the data
func (r *CommandRunner) GroupBy(field string) error {
	if r.table == nil {
		return fmt.Errorf("no data loaded")
	}
	cmd := NewGroupByCommand(r.table, field)
	table, err := cmd.Execute()
	if err != nil {
		return err
	}
	r.table = table
	return nil
}

// Output outputs the data in the specified format
func (r *CommandRunner) Output(format string) (string, error) {
	if r.table == nil {
		return "", fmt.Errorf("no data loaded")
	}
	cmd := NewFormatCommand(r.table, format)
	return cmd.Execute()
}

// GetTable returns the current table
func (r *CommandRunner) GetTable() *data.Table {
	return r.table
}

// Example demonstrates data command usage
func Example() {
	// Create sample JSON data
	jsonData := `[
		{"name": "Alice", "age": 30, "city": "New York", "status": "active"},
		{"name": "Bob", "age": 25, "city": "San Francisco", "status": "active"},
		{"name": "Charlie", "age": 35, "city": "New York", "status": "inactive"}
	]`
	
	// Create command runner
	runner := NewCommandRunner()
	
	// Parse JSON data
	if err := runner.Parse("json", strings.NewReader(jsonData)); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing: %v\n", err)
		return
	}
	
	// Filter for active users
	if err := runner.Where("status", "==", "active"); err != nil {
		fmt.Fprintf(os.Stderr, "Error filtering: %v\n", err)
		return
	}
	
	// Select specific columns
	if err := runner.Select("name", "city"); err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting: %v\n", err)
		return
	}
	
	// Output as table
	output, err := runner.Output("table")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting: %v\n", err)
		return
	}
	
	fmt.Println(output)
}
