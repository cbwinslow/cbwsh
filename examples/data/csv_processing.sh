#!/bin/bash
# Example 2: Processing CSV data with cbwsh data commands

# Sample server metrics CSV
cat > /tmp/servers.csv << 'EOF'
hostname,region,cpu_percent,memory_gb,disk_percent,status
web-01,us-east-1,45.2,8.5,62,healthy
web-02,us-east-1,78.9,15.2,85,warning
db-01,us-west-2,32.1,32.0,45,healthy
db-02,us-west-2,89.5,28.7,92,critical
api-01,eu-west-1,55.3,12.1,70,healthy
api-02,eu-west-1,23.4,6.8,35,healthy
cache-01,ap-south-1,67.8,4.2,55,warning
EOF

echo "=== Example 1: Find servers with high CPU usage ==="
# cbwsh data parse csv < /tmp/servers.csv | where cpu_percent > 70 | select hostname cpu_percent status
# Expected: Servers with CPU > 70%

echo ""
echo "=== Example 2: Group by region ==="
# cbwsh data parse csv < /tmp/servers.csv | group-by region
# Expected: Count of servers per region

echo ""
echo "=== Example 3: Critical alerts ==="
# cbwsh data parse csv < /tmp/servers.csv | where status == critical | select hostname cpu_percent memory_gb disk_percent
# Expected: Critical servers with their metrics

echo ""
echo "=== Example 4: Top 3 servers by memory usage ==="
# cbwsh data parse csv < /tmp/servers.csv | sort memory_gb desc | limit 3 | select hostname memory_gb
# Expected: Top 3 servers by memory

echo ""
echo "=== Example 5: Healthy servers in US regions ==="
# cbwsh data parse csv < /tmp/servers.csv | where status == healthy | where region contains us | select hostname region
# Expected: Healthy servers in US regions

# Cleanup
rm -f /tmp/servers.csv

echo ""
echo "Note: Actual command execution requires cbwsh-data binary"
