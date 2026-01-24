# cbwsh Data Processing Examples

This directory contains examples demonstrating the structured data processing capabilities of cbwsh, inspired by modern shells like Nushell and PowerShell.

## Overview

cbwsh's data processing features enable you to work with structured data formats (JSON, YAML, CSV) using type-safe pipelines and SQL-like operations.

## Available Examples

### 1. JSON Processing (`json_processing.sh`)

Demonstrates basic JSON data operations:
- Parsing JSON data
- Filtering records with `where`
- Selecting columns with `select`
- Sorting data
- Grouping data
- Output formatting

**Run:**
```bash
./json_processing.sh
```

### 2. CSV Processing (`csv_processing.sh`)

Shows CSV data analysis:
- Parsing CSV files
- Server metrics analysis
- Resource monitoring
- Alert detection

**Run:**
```bash
./csv_processing.sh
```

### 3. DevOps Use Cases (`devops_usecases.sh`)

Real-world DevOps scenarios:
- Docker container analysis
- Kubernetes pod monitoring
- AWS EC2 instance management
- Log file analysis
- GitHub API data processing
- System metrics monitoring
- Database query result processing

**Run:**
```bash
./devops_usecases.sh
```

## Command Syntax

### Basic Pipeline

```bash
# Parse -> Filter -> Select -> Output
cat data.json | cbwsh data parse json | where status == active | select name email
```

### Available Commands

#### Parse
Parse structured data from various formats:
```bash
cbwsh data parse json < file.json
cbwsh data parse yaml < file.yaml
cbwsh data parse csv < file.csv
```

#### Where
Filter records based on conditions:
```bash
where field == value          # Equality
where field != value          # Inequality
where field contains value    # String contains
where field > value           # Greater than (numeric)
where field < value           # Less than (numeric)
```

#### Select
Choose specific columns:
```bash
select col1 col2 col3
select name email             # Multiple columns
select *                      # All columns (default)
```

#### Sort
Sort data by a field:
```bash
sort field asc                # Ascending (default)
sort field desc               # Descending
```

#### Limit
Limit number of rows:
```bash
limit 10                      # First 10 rows
limit 100                     # First 100 rows
```

#### Group By
Group records by a field:
```bash
group-by field                # Count per unique value
```

#### To
Convert output format:
```bash
to table                      # Table format (default)
to json                       # JSON output
to csv                        # CSV output
to yaml                       # YAML output
```

## Example Workflows

### Workflow 1: Filter and Transform

```bash
# Start with JSON data
cat users.json | \
  # Parse as JSON
  cbwsh data parse json | \
  # Filter active users
  where status == active | \
  # Only show specific fields
  select name email role | \
  # Sort by name
  sort name asc | \
  # Show as table
  to table
```

### Workflow 2: Aggregation

```bash
# Server metrics
cat servers.csv | \
  cbwsh data parse csv | \
  # Group by region
  group-by region | \
  # Sort by count
  sort count desc
```

### Workflow 3: Complex Filtering

```bash
# Multi-condition filter
cat logs.json | \
  cbwsh data parse json | \
  # Error logs only
  where level == error | \
  # From last hour
  where timestamp > "2024-01-15T20:00:00Z" | \
  # Select relevant fields
  select timestamp service message | \
  # Sort by time
  sort timestamp desc | \
  # Top 20
  limit 20
```

## Integration with Standard Tools

### With curl (API data)

```bash
curl -s https://api.github.com/users/octocat/repos | \
  cbwsh data parse json | \
  where language == JavaScript | \
  select name stars forks | \
  sort stars desc | \
  limit 5
```

### With Docker

```bash
docker ps --format '{{json .}}' | \
  cbwsh data parse json --stream | \
  where status contains Up | \
  select names image ports
```

### With kubectl

```bash
kubectl get pods -o json | \
  jq '.items' | \
  cbwsh data parse json | \
  where status != Running | \
  select name namespace status
```

### With AWS CLI

```bash
aws ec2 describe-instances --output json | \
  jq '.Reservations[].Instances[]' | \
  cbwsh data parse json | \
  where state == running | \
  select instance_id instance_type
```

## Performance Tips

1. **Use streaming for large files:**
   ```bash
   cbwsh data parse json --stream < large_file.json
   ```

2. **Filter early in pipeline:**
   ```bash
   # Good: Filter first, then process
   parse json | where status == active | select name email
   
   # Less efficient: Select all, then filter
   parse json | select * | where status == active
   ```

3. **Limit output size:**
   ```bash
   parse json | limit 1000 | where condition | ...
   ```

## Error Handling

```bash
# Check if file exists
if [ -f data.json ]; then
  cat data.json | cbwsh data parse json || echo "Parse error"
else
  echo "File not found"
fi

# With error output
cat data.json | cbwsh data parse json 2>&1 | tee /tmp/errors.log
```

## Common Patterns

### Pattern 1: Top N Analysis
```bash
# Top 10 users by activity
cat users.json | parse json | sort activity desc | limit 10
```

### Pattern 2: Deduplication
```bash
# Unique values in a field
cat data.json | parse json | group-by field | select field
```

### Pattern 3: Conditional Selection
```bash
# Select based on multiple conditions
cat data.json | parse json | \
  where status == active | \
  where age > 18 | \
  where country == US
```

### Pattern 4: Format Conversion
```bash
# CSV to JSON
cat data.csv | cbwsh data parse csv | to json > data.json

# JSON to CSV
cat data.json | cbwsh data parse json | to csv > data.csv
```

## Extending Examples

To create your own examples:

1. **Start with sample data:**
   ```bash
   echo '[{"name":"test","value":42}]' | cbwsh data parse json
   ```

2. **Add filters and transforms:**
   ```bash
   echo '[...]' | parse json | where value > 40 | select name
   ```

3. **Test output formats:**
   ```bash
   echo '[...]' | parse json | to table
   echo '[...]' | parse json | to json
   ```

## Resources

- [SHELL_RESEARCH.md](../../SHELL_RESEARCH.md) - Research on modern shell features
- [SHELL_VARIANTS.md](../../SHELL_VARIANTS.md) - Shell variant documentation
- [pkg/data](../../pkg/data/) - Data processing implementation
- [pkg/commands/data](../../pkg/commands/data/) - Command implementations

## Contributing

To add more examples:

1. Create a new `.sh` file in this directory
2. Add descriptive comments
3. Include sample data
4. Show expected output
5. Update this README

## Support

For issues or questions about data processing features:

- Open an issue on GitHub
- Check the documentation
- Review existing examples
- Ask in discussions

---

*Examples last updated: January 2026*
