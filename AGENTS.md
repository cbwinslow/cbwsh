# AI Agents Configuration Guide

This guide explains how to configure and use various AI agents with cbwsh for code review, assistance, and automation in GitHub Actions and other workflows.

## Table of Contents

- [Overview](#overview)
- [CodeRabbit](#coderabbit)
- [GitHub Copilot](#github-copilot)
- [OpenAI Codex](#openai-codex)
- [Google Gemini](#google-gemini)
- [Anthropic Claude](#anthropic-claude)
- [oh-my-opencode](#oh-my-opencode)
- [Local Ollama](#local-ollama)
- [GitHub Actions Integration](#github-actions-integration)
- [Best Practices](#best-practices)

## Overview

cbwsh supports multiple AI agents for different purposes:

| Agent | Primary Use | Integration |
|-------|-------------|-------------|
| **CodeRabbit** | Automated PR reviews | GitHub App + CLI |
| **GitHub Copilot** | Code suggestions | CLI + VSCode |
| **OpenAI Codex** | Code generation | API |
| **Google Gemini** | Multimodal AI | API |
| **Anthropic Claude** | Long-context analysis | API |
| **Ollama** | Local LLM inference | Local API |

## CodeRabbit

CodeRabbit provides automated, AI-powered code reviews for pull requests.

### Setup

1. **Install CodeRabbit App** on your repository:
   - Visit [CodeRabbit App](https://coderabbit.ai/)
   - Click "Install" and select your repository
   - Grant necessary permissions

2. **Configure CodeRabbit** with `.coderabbit.yaml`:

```yaml
# .coderabbit.yaml
language: "en-US"
early_access: false
reviews:
  profile: "chill"
  request_changes_workflow: false
  high_level_summary: true
  poem: true
  review_status: true
  collapse_walkthrough: false
  auto_review:
    enabled: true
    drafts: false
  tools:
    github-checks:
      enabled: true
      timeout: 90
    ruff:
      enabled: true
    shellcheck:
      enabled: true
    yamllint:
      enabled: true
    golangci-lint:
      enabled: true
chat:
  auto_reply: true
```

### CodeRabbit CLI

Install the CodeRabbit CLI for local reviews:

```bash
# Install via npm
npm install -g @coderabbitai/cli

# Or via brew (macOS)
brew install coderabbit/tap/coderabbit

# Login
coderabbit auth login

# Review current changes
coderabbit review

# Review specific files
coderabbit review src/*.go
```

### Usage in cbwsh

```bash
# From within cbwsh
coderabbit review
coderabbit chat "How can I optimize this function?"
```

## GitHub Copilot

GitHub Copilot provides AI-powered code suggestions.

### Setup

1. **Install GitHub Copilot CLI**:

```bash
# Install via npm
npm install -g @githubnext/github-copilot-cli

# Setup authentication
github-copilot-cli auth

# Add aliases (optional)
eval "$(github-copilot-cli alias -- "$0")"
```

2. **Configure in cbwsh**:

Add to your `~/.cbwsh/config.yaml`:

```yaml
ai:
  copilot:
    enabled: true
    api_key: ""  # Set via GitHub auth
```

### Usage

```bash
# In cbwsh, ask Copilot for command help
?? how to find large files
# Suggests: find . -type f -size +100M

# Explain a command
?! tar -xzvf archive.tar.gz
# Explains what the command does

# Get Git help
git? undo last commit
# Suggests: git reset --soft HEAD~1
```

## OpenAI Codex

OpenAI's Codex powers code generation and understanding.

### Setup

1. **Get API Key**:
   - Visit [OpenAI Platform](https://platform.openai.com/)
   - Create an API key
   - Store securely in cbwsh secrets

2. **Configure**:

```yaml
# ~/.cbwsh/config.yaml
ai:
  provider: openai
  api_key: ${OPENAI_API_KEY}  # Or store in secrets
  model: gpt-4
  max_tokens: 2000
  temperature: 0.2
```

3. **Store API Key**:

```bash
# In cbwsh, store the key securely
cbwsh secrets set OPENAI_API_KEY "sk-..."
```

### Usage in cbwsh

```bash
# Press Ctrl+A to enter AI mode
# Then ask:
"Generate a bash script to backup my home directory"
"Explain this error: permission denied"
"How do I parse JSON in bash?"
```

## Google Gemini

Google's Gemini provides multimodal AI capabilities.

### Setup

1. **Get API Key**:
   - Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
   - Create an API key

2. **Configure**:

```yaml
# ~/.cbwsh/config.yaml
ai:
  provider: gemini
  api_key: ${GEMINI_API_KEY}
  model: gemini-pro
  temperature: 0.3
```

### Usage

```bash
# Gemini excels at code analysis and generation
# Press Ctrl+A in cbwsh:
"Analyze this function for security issues"
"Generate unit tests for my script"
"Optimize this loop for performance"
```

## Anthropic Claude

Claude by Anthropic excels at long-context analysis and reasoning.

### Setup

1. **Get API Key**:
   - Visit [Anthropic Console](https://console.anthropic.com/)
   - Create an API key

2. **Configure**:

```yaml
# ~/.cbwsh/config.yaml
ai:
  provider: anthropic
  api_key: ${ANTHROPIC_API_KEY}
  model: claude-3-opus-20240229
  max_tokens: 4096
```

### Usage

```bash
# Claude is excellent for complex analysis
# Press Ctrl+A:
"Review this entire codebase structure"
"Explain the architecture of this system"
"Debug this complex error trace"
```

## oh-my-opencode

oh-my-opencode is a framework for integrating AI into development workflows.

### Setup

1. **Install**:

```bash
# Clone the repository
git clone https://github.com/oh-my-opencode/oh-my-opencode.git
cd oh-my-opencode

# Install dependencies
npm install

# Setup
npm run setup
```

2. **Configure**:

Create `.oh-my-opencode.yaml`:

```yaml
version: "1.0"
ai:
  providers:
    - name: openai
      api_key: ${OPENAI_API_KEY}
    - name: anthropic
      api_key: ${ANTHROPIC_API_KEY}
  
agents:
  - name: code-reviewer
    provider: openai
    model: gpt-4
    temperature: 0.2
    
  - name: documentation-writer
    provider: anthropic
    model: claude-3-opus
    
  - name: test-generator
    provider: openai
    model: gpt-4

workflows:
  pr-review:
    trigger: pull_request
    agents: [code-reviewer]
    
  docs-update:
    trigger: push
    agents: [documentation-writer]
```

### Integration with cbwsh

```bash
# Use oh-my-opencode agents from cbwsh
omc review --agent code-reviewer
omc generate-docs --agent documentation-writer
omc test --agent test-generator
```

## Local Ollama

Ollama runs AI models locally for privacy and offline use.

### Setup

1. **Install Ollama**:

```bash
# Linux/macOS
curl https://ollama.ai/install.sh | sh

# Or with Homebrew
brew install ollama

# Windows (via WSL)
curl https://ollama.ai/install.sh | sh
```

2. **Pull Models**:

```bash
# Download models
ollama pull llama2
ollama pull codellama
ollama pull mistral

# List available models
ollama list
```

3. **Configure in cbwsh**:

```yaml
# ~/.cbwsh/config.yaml
ai:
  provider: ollama
  ollama_url: http://localhost:11434
  ollama_model: codellama
  enable_monitoring: true
  monitoring_interval: 30
```

### Usage

```bash
# Start Ollama service
ollama serve

# In cbwsh, AI features will use Ollama
# Press Ctrl+M to enable AI monitor
# Press Ctrl+A for AI assist

# Use specific models
cbwsh ai-chat --model llama2
cbwsh ai-suggest --model codellama
```

### Available Models

| Model | Size | Best For |
|-------|------|----------|
| **llama2** | 7B/13B/70B | General purpose |
| **codellama** | 7B/13B/34B | Code generation |
| **mistral** | 7B | Fast inference |
| **mixtral** | 8x7B | High quality |
| **deepseek-coder** | 6.7B | Code completion |

## GitHub Actions Integration

Integrate AI agents into your CI/CD workflows.

### CodeRabbit Action

Create `.github/workflows/coderabbit.yml`:

```yaml
name: CodeRabbit Review

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: CodeRabbit Review
        uses: coderabbitai/coderabbit-action@v1
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
```

### OpenAI Code Review

Create `.github/workflows/ai-review.yml`:

```yaml
name: AI Code Review

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  ai-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
      
      - name: Install Dependencies
        run: |
          npm install -g openai-cli
      
      - name: AI Code Review
        env:
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
        run: |
          # Get diff
          git diff origin/${{ github.base_ref }}...HEAD > changes.diff
          
          # Run AI review
          openai-cli review changes.diff \
            --model gpt-4 \
            --output review.md
          
      - name: Post Review Comment
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const review = fs.readFileSync('review.md', 'utf8');
            
            github.rest.issues.createComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              body: `## 🤖 AI Code Review\n\n${review}`
            });
```

### Gemini Analysis

Create `.github/workflows/gemini-analysis.yml`:

```yaml
name: Gemini Code Analysis

on:
  push:
    branches: [main, develop]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      
      - name: Install Gemini SDK
        run: |
          pip install google-generativeai
      
      - name: Run Analysis
        env:
          GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
        run: |
          python << 'EOF'
          import google.generativeai as genai
          import os
          
          genai.configure(api_key=os.environ['GEMINI_API_KEY'])
          model = genai.GenerativeModel('gemini-pro')
          
          # Analyze code
          with open('main.go', 'r') as f:
              code = f.read()
          
          prompt = f"Analyze this Go code for security issues:\n\n{code}"
          response = model.generate_content(prompt)
          
          print(response.text)
          EOF
```

### Multi-Agent Workflow

Create `.github/workflows/multi-agent.yml`:

```yaml
name: Multi-Agent Review

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  coderabbit-review:
    name: CodeRabbit Review
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: coderabbitai/coderabbit-action@v1
  
  copilot-suggestions:
    name: Copilot Suggestions
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Copilot CLI
        run: |
          npm install -g @githubnext/github-copilot-cli
          # Add suggestions logic
  
  gemini-security:
    name: Gemini Security Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Security Analysis
        env:
          GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
        run: |
          # Run security analysis with Gemini
          python scripts/security-check.py
  
  claude-architecture:
    name: Claude Architecture Review
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Architecture Analysis
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          # Analyze architecture with Claude
          python scripts/architecture-review.py
```

## Best Practices

### Security

1. **Never commit API keys** - Use GitHub Secrets
2. **Rotate keys regularly** - Update secrets every 90 days
3. **Use minimal permissions** - Grant only necessary access
4. **Store keys in cbwsh secrets**:
   ```bash
   cbwsh secrets set OPENAI_API_KEY "sk-..."
   cbwsh secrets set ANTHROPIC_API_KEY "sk-ant-..."
   ```

### Cost Management

1. **Set token limits** in AI configurations
2. **Use local models** (Ollama) for development
3. **Cache responses** when possible
4. **Monitor usage**:
   ```yaml
   ai:
     max_tokens: 1000
     enable_caching: true
     usage_tracking: true
   ```

### Performance

1. **Use appropriate models** for tasks:
   - Simple queries: Smaller models (llama2-7B)
   - Complex analysis: Larger models (gpt-4, claude-3-opus)
   
2. **Enable streaming** for faster responses:
   ```yaml
   ai:
     streaming: true
   ```

3. **Batch similar requests**

### Privacy

1. **Use local models** for sensitive code:
   ```yaml
   ai:
     provider: ollama
     ollama_model: codellama
   ```

2. **Review data sharing policies** of AI providers

3. **Consider self-hosted options**:
   - Ollama for local inference
   - LocalAI for OpenAI-compatible API
   - Text-generation-webui for advanced local setups

## Troubleshooting

### Common Issues

**API Key Errors**:
```bash
# Verify key is set
cbwsh secrets list

# Test connection
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

**Ollama Connection Issues**:
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Restart Ollama
ollama serve
```

**Rate Limiting**:
```yaml
# Add rate limiting configuration
ai:
  rate_limit:
    requests_per_minute: 20
    retry_on_limit: true
```

## Resources

- **CodeRabbit**: https://coderabbit.ai/
- **GitHub Copilot**: https://github.com/features/copilot
- **OpenAI Platform**: https://platform.openai.com/
- **Google AI Studio**: https://makersuite.google.com/
- **Anthropic Console**: https://console.anthropic.com/
- **Ollama**: https://ollama.ai/
- **oh-my-opencode**: https://github.com/oh-my-opencode/oh-my-opencode

---

For more information, see [cbwsh documentation](https://github.com/cbwinslow/cbwsh) or open an issue.
