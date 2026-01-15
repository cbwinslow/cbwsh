#!/bin/bash
# Example 3: Real-world DevOps use cases with cbwsh data commands

echo "=== Use Case 1: Docker Container Analysis ==="
echo "Parse docker ps output and find containers using too much memory"
cat > /tmp/docker_demo.sh << 'EOF'
# Get docker containers as JSON
docker ps --format '{{json .}}' | \
  # Parse JSON and filter high memory
  cbwsh data parse json --stream | \
  where memory_percent > 80 | \
  select name image memory_percent status | \
  sort memory_percent desc
EOF
echo "Script created: /tmp/docker_demo.sh"
echo ""

echo "=== Use Case 2: Kubernetes Pod Analysis ==="
echo "Find pods in non-running state"
cat > /tmp/k8s_demo.sh << 'EOF'
# Get pods as JSON
kubectl get pods -o json | \
  # Extract items array
  jq '.items' | \
  # Parse and filter
  cbwsh data parse json | \
  where status != Running | \
  select name namespace status restart_count | \
  sort restart_count desc
EOF
echo "Script created: /tmp/k8s_demo.sh"
echo ""

echo "=== Use Case 3: AWS EC2 Instance Analysis ==="
echo "Find instances by tag and state"
cat > /tmp/aws_demo.sh << 'EOF'
# Get EC2 instances
aws ec2 describe-instances --output json | \
  # Parse and filter
  cbwsh data parse json | \
  where state == running | \
  where environment == production | \
  select instance_id instance_type availability_zone | \
  group-by instance_type
EOF
echo "Script created: /tmp/aws_demo.sh"
echo ""

echo "=== Use Case 4: Log Analysis ==="
echo "Parse nginx access logs in JSON format"
cat > /tmp/log_demo.sh << 'EOF'
# Parse nginx JSON logs
cat /var/log/nginx/access.log.json | \
  cbwsh data parse json --stream | \
  where status >= 400 | \
  select timestamp remote_addr request status response_time | \
  sort response_time desc | \
  limit 20
EOF
echo "Script created: /tmp/log_demo.sh"
echo ""

echo "=== Use Case 5: GitHub Repository Analysis ==="
echo "Analyze repository statistics"
cat > /tmp/github_demo.sh << 'EOF'
# Get repos from GitHub API
curl -s "https://api.github.com/orgs/kubernetes/repos?per_page=100" | \
  cbwsh data parse json | \
  where language == Go | \
  select name stargazers_count forks_count open_issues_count | \
  sort stargazers_count desc | \
  limit 10
EOF
echo "Script created: /tmp/github_demo.sh"
echo ""

echo "=== Use Case 6: System Monitoring ==="
echo "Parse system metrics"
cat > /tmp/system_demo.sh << 'EOF'
# Parse custom metrics in CSV
cat /var/log/metrics.csv | \
  cbwsh data parse csv | \
  where cpu_usage > 80 or memory_usage > 90 | \
  select timestamp hostname cpu_usage memory_usage | \
  sort timestamp desc
EOF
echo "Script created: /tmp/system_demo.sh"
echo ""

echo "=== Use Case 7: Database Query Results ==="
echo "Process PostgreSQL query results"
cat > /tmp/db_demo.sh << 'EOF'
# Export query results to JSON and analyze
psql -t -A -F, -c "SELECT * FROM users" | \
  cbwsh data parse csv --no-header | \
  where subscription_status == active | \
  where created_at > "2024-01-01" | \
  group-by country | \
  sort count desc
EOF
echo "Script created: /tmp/db_demo.sh"
echo ""

echo "All demo scripts created in /tmp/"
echo "These demonstrate real-world use cases for cbwsh data processing"
echo ""
echo "Note: Actual execution requires:"
echo "  1. cbwsh-data binary installed"
echo "  2. Appropriate services/tools available (docker, kubectl, aws, etc.)"
echo "  3. Valid credentials and permissions"
