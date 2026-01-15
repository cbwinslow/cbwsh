# cbwsh Shell Variants - Implementation Guide

This document describes the shell variants available in cbwsh, based on GitHub research of top shell projects.

## Overview

Based on analysis of 100+ top GitHub shell repositories, cbwsh implements several shell variants optimized for different use cases. Each variant shares the core cbwsh functionality while adding specialized features.

## Shell Variants

### 1. cbwsh (Standard)

**Purpose:** General-purpose modern shell
**Status:** ✅ Implemented

**Features:**
- Multi-shell support (bash, zsh)
- AI integration (OpenAI, Anthropic, Gemini, Ollama)
- Autocompletion and syntax highlighting
- Multi-pane support
- SSH management
- Secrets management
- Plugin system

**Target Users:** General users, developers

**Usage:**
```bash
cbwsh
```

### 2. cbwsh-data (Data Processing Variant)

**Purpose:** Shell optimized for data processing and analysis
**Status:** 🚧 In Progress
**Inspired by:** Nushell, PowerShell

**Key Features:**
- Structured data pipelines (type-safe)
- Native JSON/YAML/CSV parsing
- SQL-like query operations
- Table visualization
- Data transformation commands
- Format conversion utilities

**New Commands:**
```bash
# Parse JSON and query
cat data.json | cbwsh-data parse json | where status == "active" | select name, email

# CSV analysis
cat users.csv | cbwsh-data parse csv | group-by country | sort count desc

# YAML transformation
cat config.yaml | cbwsh-data parse yaml | select service.* | to json
```

**Target Users:** Data engineers, analysts, DevOps engineers

**Implementation Status:**
- ✅ Core data types (pkg/data/types.go)
- ✅ JSON/YAML/CSV parsers (pkg/data/parsers.go)
- ✅ Pipeline operations (pkg/data/pipeline.go)
- ⏳ Command implementations
- ⏳ Integration with main shell

### 3. cbwsh-ai (AI-Powered Variant)

**Purpose:** AI-first shell for maximum productivity
**Status:** ⏳ Planned
**Inspired by:** Warp, Shell GPT, Amazon Q Developer CLI

**Key Features:**
- Natural language command generation
- Context-aware AI assistant
- Multi-agent workflows
- Code generation and review
- Learning from user patterns
- Proactive suggestions
- Error analysis and auto-fix

**Example Interactions:**
```bash
# Natural language commands
cbwsh-ai> find all large log files older than 7 days
# Generates: find /var/log -name "*.log" -size +100M -mtime +7

# Error auto-fix
$ rm -rf /important/file
Error: Permission denied
cbwsh-ai: Detected permission error. Suggested fix:
  sudo rm -rf /important/file
Apply? [y/n]

# Code generation
cbwsh-ai> create a bash script to backup mysql database
# Generates complete script with error handling
```

**Target Users:** Developers, power users, AI enthusiasts

### 4. cbwsh-security (Security Variant)

**Purpose:** Security-focused shell with enhanced protections
**Status:** ⏳ Planned
**Inspired by:** lynis, PEASS-ng

**Key Features:**
- Enhanced secrets management
- Security auditing tools
- Compliance checking
- Command sandboxing
- Audit logging
- Intrusion detection integration
- Secure command execution

**Security Features:**
```bash
# Security audit
cbwsh-security audit system

# Secrets scanning
cbwsh-security scan secrets .

# Compliance check
cbwsh-security check compliance pci-dss

# Secure execution
cbwsh-security exec --sandbox "untrusted-script.sh"
```

**Target Users:** Security professionals, system administrators

### 5. cbwsh-devops (DevOps Variant)

**Purpose:** DevOps-optimized shell with cloud and container tools
**Status:** ⏳ Planned
**Inspired by:** k9s, lazydocker, aws-cli

**Key Features:**
- Docker/Kubernetes UI integration
- Cloud CLI helpers (AWS, GCP, Azure)
- Infrastructure monitoring
- Log analysis tools
- CI/CD integration
- Terraform/Ansible helpers

**DevOps Commands:**
```bash
# Docker management
cbwsh-devops docker ps --format table

# Kubernetes operations
cbwsh-devops k8s pods --namespace production | where status != "Running"

# AWS resources
cbwsh-devops aws ec2 list | where region == "us-east-1" | select id, type, state

# Log analysis
cbwsh-devops logs analyze /var/log/app.log --errors
```

**Target Users:** DevOps engineers, SREs, cloud architects

## Technical Architecture

### Shared Core

All variants share the core cbwsh components:

```
pkg/
├── core/           # Core types and interfaces
├── config/         # Configuration management
├── shell/          # Shell executors
├── panes/          # Pane management
├── plugins/        # Plugin system
├── secrets/        # Secrets management
├── ssh/            # SSH support
├── ai/             # AI integration
└── ui/             # UI components
```

### Variant-Specific Components

Each variant adds specialized packages:

**cbwsh-data:**
```
pkg/data/
├── types.go        # Data types (Table, Record, Value)
├── parsers.go      # Format parsers (JSON, YAML, CSV)
├── pipeline.go     # Pipeline operations
└── formatters.go   # Output formatters
```

**cbwsh-ai:**
```
pkg/ai/
├── agents/         # Specialized AI agents
├── nlp/           # Natural language processing
├── learning/      # Pattern learning
└── context/       # Context management
```

**cbwsh-security:**
```
pkg/security/
├── audit/         # Security auditing
├── compliance/    # Compliance checking
├── sandbox/       # Command sandboxing
└── scanner/       # Vulnerability scanning
```

**cbwsh-devops:**
```
pkg/devops/
├── docker/        # Docker integration
├── k8s/          # Kubernetes integration
├── cloud/        # Cloud provider helpers
└── monitoring/   # Infrastructure monitoring
```

## Implementation Patterns

### Pattern 1: Plugin-Based Extensions

Variants use plugins to add functionality:

```go
// Example: Data processing plugin
type DataPlugin struct {
    parser data.Parser
}

func (p *DataPlugin) Execute(args []string) error {
    // Parse and process data
    table, err := p.parser.Parse(input)
    if err != nil {
        return err
    }
    
    // Apply pipeline operations
    result := data.NewPipeline(table).
        Where("status", func(v *data.Value) bool {
            return v.String() == "active"
        }).
        Select("name", "email").
        Execute()
    
    // Format and output
    fmt.Println(formatters.NewTableFormatter().Format(result))
    return nil
}
```

### Pattern 2: Command Wrappers

Enhance existing commands with structured output:

```go
// Wrap 'docker ps' with structured output
func WrapDockerPS() *data.Table {
    output := exec.Command("docker", "ps", "--format", "json").Output()
    parser := data.NewJSONParser()
    return parser.Parse(output)
}
```

### Pattern 3: AI Integration

Use AI for command enhancement:

```go
// Natural language to command
func TranslateNaturalLanguage(prompt string) (string, error) {
    agent := ai.NewAgent("translator", core.AIProviderOllama)
    return agent.TranslateCommand(context.Background(), prompt)
}
```

## Building Variants

### From Source

```bash
# Build standard cbwsh
make build

# Build specific variant
make build-data
make build-ai
make build-security
make build-devops

# Build all variants
make build-all
```

### Using Go

```bash
# Install standard
go install github.com/cbwinslow/cbwsh@latest

# Install variants
go install github.com/cbwinslow/cbwsh/cmd/cbwsh-data@latest
go install github.com/cbwinslow/cbwsh/cmd/cbwsh-ai@latest
```

## Configuration

Each variant has its own configuration section in `~/.cbwsh/config.yaml`:

```yaml
# Standard cbwsh config
shell:
  default_shell: bash
  
ui:
  theme: default

# Data variant config
data:
  default_format: json
  max_table_rows: 1000
  pretty_print: true

# AI variant config
ai:
  provider: ollama
  model: llama2
  context_window: 4096
  learning_enabled: true

# Security variant config
security:
  audit_log: ~/.cbwsh/audit.log
  strict_mode: true
  sandbox_enabled: true

# DevOps variant config
devops:
  cloud_providers:
    - aws
    - gcp
    - azure
  kubernetes_contexts:
    - production
    - staging
```

## Switching Between Variants

Users can easily switch between variants:

```bash
# Use standard cbwsh
cbwsh

# Use data variant
cbwsh --variant data
# or
cbwsh-data

# Use AI variant
cbwsh --variant ai
# or
cbwsh-ai

# Use multiple variants in panes
cbwsh --left data --right ai
```

## Migration Path

Users can migrate from standard cbwsh to variants incrementally:

1. **Phase 1:** Use standard cbwsh with new data commands
2. **Phase 2:** Enable variant features via config
3. **Phase 3:** Switch to dedicated variant binary
4. **Phase 4:** Customize variant-specific features

## Research-Driven Features

Based on analysis of top GitHub shell projects, the following features are prioritized:

### High Impact (Implemented or In Progress)
- ✅ Syntax highlighting (Fish pattern)
- ✅ Autosuggestions (Fish pattern)
- ✅ AI integration (Warp, Shell GPT pattern)
- 🚧 Structured data pipelines (Nushell pattern)
- ✅ Multi-pane support (tmux pattern)
- ✅ Plugin system (Zsh pattern)

### Medium Impact (Planned)
- ⏳ Natural language commands (Warp pattern)
- ⏳ Container integration (lazydocker pattern)
- ⏳ Cloud CLI helpers (aws-cli pattern)
- ⏳ Security auditing (lynis pattern)

### Future Enhancements
- Context-aware completions
- Machine learning suggestions
- Workflow automation
- Team collaboration features

## Contributing

To contribute a new shell variant:

1. Create variant-specific packages in `pkg/`
2. Implement variant commands
3. Add configuration schema
4. Write tests
5. Update documentation
6. Submit PR with research justification

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## Performance Considerations

Variants are designed to be lightweight:

- **Startup time:** <100ms target for all variants
- **Memory:** Minimal overhead per variant
- **Lazy loading:** Features loaded on-demand
- **Plugin system:** Optional features don't affect core performance

## Compatibility

All variants maintain compatibility with:

- Standard shell scripts (bash, zsh)
- Existing tools and commands
- Plugin ecosystem
- Configuration files

## Future Roadmap

See [ROADMAP.md](ROADMAP.md) for detailed plans on:

- Additional variants
- Feature enhancements
- Performance improvements
- Integration with emerging tools

---

*Last updated: January 2026*
*Based on research of 100+ top GitHub shell repositories*
