# Integration Guide for cbwsh

## Overview

This guide helps you integrate cbwsh into your existing workflow, especially if you're coming from other shells like Bash, Zsh, Fish, or modern alternatives like Nushell.

## Quick Migration

### From Bash/Zsh

cbwsh is fully compatible with Bash and Zsh commands. Your existing scripts and workflows will work without modification.

```bash
# Your existing commands work as-is
ls -la
cd ~/projects
git status
docker ps

# Plus you get cbwsh enhancements
ls -la | cbwsh data parse json  # Structured data processing
```

### From Fish

If you're used to Fish's friendly features:

| Fish Feature | cbwsh Equivalent | Status |
|-------------|------------------|--------|
| Autosuggestions | ✅ Built-in | Implemented |
| Syntax highlighting | ✅ Built-in | Implemented |
| Web-based config | ⏳ Planned | Future |
| Abbreviations | ⏳ Planned | Future |

### From Nushell

cbwsh's data processing is inspired by Nushell:

| Nushell Feature | cbwsh Equivalent | Example |
|----------------|------------------|---------|
| `open file.json` | `cat file.json \| cbwsh data parse json` | Parse JSON |
| `where status == "active"` | `where status == active` | Filter data |
| `select name email` | `select name email` | Select columns |
| `sort-by age` | `sort age desc` | Sort data |
| `group-by region` | `group-by region` | Group data |

## Installation Alongside Existing Shells

### Keep Your Current Shell

You don't have to replace your shell. Use cbwsh alongside:

```bash
# In ~/.bashrc or ~/.zshrc
alias cbw='cbwsh'

# Quick data processing without switching shells
alias jq-table='cbwsh data parse json | to table'
alias csv-view='cbwsh data parse csv | to table'
```

### Gradual Adoption

**Phase 1: Try it out**
```bash
# Run cbwsh when you need it
cbwsh
# Exit back to your shell: Ctrl+Q
```

**Phase 2: Use for specific tasks**
```bash
# Process data with cbwsh, pipe to your shell
cat data.json | cbwsh data parse json | where active == true | to json > filtered.json
```

**Phase 3: Make it your default**
```bash
# Set cbwsh as default shell (optional)
chsh -s $(which cbwsh)
```

## Feature Integration

### 1. Structured Data Processing

#### Replacing jq

Before (jq):
```bash
cat users.json | jq '.[] | select(.status=="active") | {name, email}'
```

After (cbwsh):
```bash
cat users.json | cbwsh data parse json | where status == active | select name email
```

#### Replacing awk/grep for CSV

Before:
```bash
cat servers.csv | awk -F',' '$3 > 80 {print $1,$3}'
```

After (cbwsh):
```bash
cat servers.csv | cbwsh data parse csv | where cpu_percent > 80 | select hostname cpu_percent
```

### 2. AI-Powered Assistance

#### Command Suggestions

Instead of searching documentation:

```bash
# Traditional way
man docker
google "docker remove all containers"

# cbwsh way (Ctrl+A for AI assist)
# Type: "remove all docker containers"
# AI suggests: docker rm $(docker ps -aq)
```

#### Error Fixes

```bash
# Traditional way
$ rm -rf /protected/file
Permission denied
$ sudo rm -rf /protected/file

# cbwsh way
$ rm -rf /protected/file
Error: Permission denied
AI suggests: sudo rm -rf /protected/file
Apply? [y/n]
```

### 3. SSH Management

#### Replacing ~/.ssh/config

Traditional:
```bash
# ~/.ssh/config
Host myserver
    HostName example.com
    User admin
    Port 2222
```

cbwsh:
```bash
# Save connection in cbwsh
cbwsh> ssh save myserver admin@example.com:2222

# Quick connect
cbwsh> ssh connect myserver
```

### 4. Secrets Management

#### Replacing pass/1password CLI

Traditional:
```bash
# Using pass
pass show api/github-token

# Using 1password CLI
op item get "GitHub Token" --fields password
```

cbwsh:
```bash
# Store encrypted secret
cbwsh> secrets set github-token <value>

# Retrieve secret
cbwsh> secrets get github-token

# Use in command
export GITHUB_TOKEN=$(cbwsh secrets get github-token)
```

## Workflow Integration

### Git Workflows

```bash
# Traditional git workflow
git status
git add .
git commit -m "Update files"
git push

# Enhanced with cbwsh
# Git status in prompt (built-in)
# AI-generated commit messages
cbwsh> git commit
AI suggests: "Add user authentication feature"
Apply? [y/n]
```

### Docker Workflows

```bash
# Traditional
docker ps
docker logs container-id
docker stats

# With cbwsh data processing
docker ps --format '{{json .}}' | \
  cbwsh data parse json | \
  where status contains Up | \
  select names image ports
```

### Kubernetes Workflows

```bash
# Traditional
kubectl get pods
kubectl describe pod mypod

# With cbwsh
kubectl get pods -o json | \
  jq '.items' | \
  cbwsh data parse json | \
  where status != Running | \
  select name namespace status restart_count | \
  sort restart_count desc
```

### Cloud CLI Integration

```bash
# AWS
aws ec2 describe-instances --output json | \
  jq '.Reservations[].Instances[]' | \
  cbwsh data parse json | \
  where state == running | \
  group-by instance_type

# GCP
gcloud compute instances list --format=json | \
  cbwsh data parse json | \
  where zone contains us-central1 | \
  select name machineType status

# Azure
az vm list --output json | \
  cbwsh data parse json | \
  where powerState == "VM running" | \
  select name resourceGroup location
```

## Configuration Integration

### Importing Existing Configs

#### From .bashrc/.zshrc

```bash
# Your aliases work in cbwsh
# cbwsh sources these automatically if using bash/zsh backend

# Or add to ~/.cbwsh/config.yaml
shell:
  source_files:
    - ~/.bashrc
    - ~/.bash_aliases
```

#### From Fish config

```bash
# Convert fish abbreviations to cbwsh aliases
# fish: abbr -a gco 'git checkout'
# cbwsh config:
aliases:
  gco: git checkout
  gst: git status
  glog: git log --oneline
```

#### From Starship

```bash
# cbwsh will support Starship prompts
# For now, configure prompt in ~/.cbwsh/config.yaml
ui:
  prompt:
    format: "[{user}@{host}] {cwd} {git_branch} > "
```

### Environment Variables

```bash
# These work automatically in cbwsh
export PATH=$PATH:~/.local/bin
export EDITOR=vim
export LANG=en_US.UTF-8

# cbwsh-specific environment
export CBWSH_AI_PROVIDER=ollama
export CBWSH_AI_MODEL=llama2
```

## Scripting with cbwsh

### Bash Scripts

Your Bash scripts work as-is:

```bash
#!/bin/bash
# existing-script.sh

# All these work in cbwsh
for file in *.txt; do
  echo "Processing $file"
  cat "$file" | some-command
done
```

### Enhanced Scripts with Data Processing

```bash
#!/bin/bash
# enhanced-script.sh

# Traditional processing
cat data.json > /tmp/data.json

# Enhanced with cbwsh data processing
cat /tmp/data.json | \
  cbwsh data parse json | \
  where status == active | \
  select id name | \
  to json > active-users.json
```

### Mixing cbwsh and Traditional Commands

```bash
#!/bin/bash
# mixed-script.sh

# Traditional text processing
grep ERROR /var/log/app.log > errors.txt

# cbwsh data processing
cat metrics.json | \
  cbwsh data parse json | \
  where value > 100 | \
  to csv > high-metrics.csv

# Back to traditional
wc -l errors.txt
```

## Terminal Integration

### tmux

```bash
# cbwsh works in tmux panes
# .tmux.conf
bind-key c new-window cbwsh
bind-key % split-window -h cbwsh
bind-key '"' split-window -v cbwsh
```

### Screen

```bash
# .screenrc
shell cbwsh
```

### SSH Connections

```bash
# Run cbwsh on remote hosts
ssh user@host -t cbwsh

# Or in .bashrc on remote:
[ -x "$(command -v cbwsh)" ] && exec cbwsh
```

## IDE Integration

### VS Code

```json
// settings.json
{
  "terminal.integrated.defaultProfile.linux": "cbwsh",
  "terminal.integrated.profiles.linux": {
    "cbwsh": {
      "path": "/usr/local/bin/cbwsh"
    }
  }
}
```

### JetBrains IDEs

```
Settings > Tools > Terminal > Shell path: /usr/local/bin/cbwsh
```

## Performance Considerations

### For Heavy Data Processing

```bash
# Use streaming for large files
cat huge.json | cbwsh data parse json --stream | process...

# Filter early to reduce data
cat data.json | cbwsh data parse json | where important == true | ...
```

### For Scripts

```bash
# For simple text processing, traditional tools might be faster
# Use cbwsh for structured data

# Good for cbwsh:
cat data.json | cbwsh data parse json | where status == active

# Traditional tools fine:
grep "pattern" file.txt
```

## Troubleshooting

### Command Not Found

```bash
# Check if cbwsh recognizes the command
which command-name

# If not, ensure PATH is correct
echo $PATH

# Source your shell config
source ~/.bashrc  # or ~/.zshrc
```

### Parsing Errors

```bash
# Verify data format
cat data.json | jq .  # Test with jq first

# Try explicit parsing
cat data.json | cbwsh data parse json --verbose
```

### Performance Issues

```bash
# Use native tools for simple tasks
# Traditional
cat file.txt | grep pattern

# Only use cbwsh for structured data
cat data.json | cbwsh data parse json | ...
```

## Best Practices

### 1. Use the Right Tool

```bash
# Good: Structured data processing
cat data.json | cbwsh data parse json | where field == value

# Overkill: Simple text processing
cat file.txt | cbwsh data parse... # Use grep instead
```

### 2. Gradual Migration

- Start with data processing features
- Add AI assistance gradually
- Migrate to full-time use when comfortable

### 3. Keep Scripts Portable

```bash
# Check if cbwsh is available
if command -v cbwsh &> /dev/null; then
  cat data.json | cbwsh data parse json | ...
else
  cat data.json | jq '...' | ...
fi
```

### 4. Document Usage

```bash
# Add comments for cbwsh-specific features
# This uses cbwsh data processing
cat users.json | \
  cbwsh data parse json | \
  where status == active | \
  select name email
```

## Getting Help

### In-Shell Help

```bash
# Press Ctrl+? for help
# Or
cbwsh help

# Command-specific help
cbwsh data help
cbwsh data parse --help
```

### Documentation

- [Usage Guide](USAGE.md) - Complete feature documentation
- [Shell Research](SHELL_RESEARCH.md) - Background on design decisions
- [Shell Variants](SHELL_VARIANTS.md) - Different shell variants
- [Examples](examples/data/) - Practical examples

### Community

- GitHub Issues: Report bugs or request features
- GitHub Discussions: Ask questions and share tips
- Examples: Check [examples/](examples/) directory

## Next Steps

1. **Try it out**: Run `cbwsh` in your terminal
2. **Explore features**: Try data processing on a JSON/CSV file
3. **Integrate gradually**: Add useful aliases to your workflow
4. **Provide feedback**: Share what works and what doesn't

---

*Last updated: January 2026*
*For questions or suggestions, open an issue on GitHub*
