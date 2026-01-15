#!/bin/bash
# Example 1: Processing JSON data with cbwsh data commands

# Sample users data
cat > /tmp/users.json << 'EOF'
[
  {"id": 1, "name": "Alice", "email": "alice@example.com", "age": 30, "status": "active", "role": "admin"},
  {"id": 2, "name": "Bob", "email": "bob@example.com", "age": 25, "status": "active", "role": "user"},
  {"id": 3, "name": "Charlie", "email": "charlie@example.com", "age": 35, "status": "inactive", "role": "user"},
  {"id": 4, "name": "Diana", "email": "diana@example.com", "age": 28, "status": "active", "role": "moderator"},
  {"id": 5, "name": "Eve", "email": "eve@example.com", "age": 32, "status": "active", "role": "admin"}
]
EOF

echo "=== Example 1: Filter active users ==="
# cbwsh data parse json < /tmp/users.json | where status == active
# Expected: Shows only active users

echo ""
echo "=== Example 2: Select specific columns ==="
# cbwsh data parse json < /tmp/users.json | select name email role
# Expected: Shows only name, email, and role columns

echo ""
echo "=== Example 3: Filter and select combined ==="
# cbwsh data parse json < /tmp/users.json | where role == admin | select name email
# Expected: Shows name and email of admin users only

echo ""
echo "=== Example 4: Sort by age ==="
# cbwsh data parse json < /tmp/users.json | sort age desc
# Expected: Users sorted by age in descending order

echo ""
echo "=== Example 5: Group by role ==="
# cbwsh data parse json < /tmp/users.json | group-by role
# Expected: Count of users per role

echo ""
echo "=== Example 6: Complex query ==="
# cbwsh data parse json < /tmp/users.json | where status == active | where age > 27 | select name age role | sort age asc
# Expected: Active users over 27, showing name/age/role, sorted by age

echo ""
echo "=== Example 7: Output as JSON ==="
# cbwsh data parse json < /tmp/users.json | where role == admin | to json
# Expected: Admin users in JSON format

# Cleanup
rm -f /tmp/users.json

echo ""
echo "Note: Actual command execution requires cbwsh-data binary"
echo "These examples demonstrate the intended usage patterns"
