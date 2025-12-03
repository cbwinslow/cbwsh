# cbwsh PowerShell Installation Script
# ======================================
# A modern, modular terminal shell built with Bubble Tea
#
# This script installs cbwsh on Windows systems
#
# Usage:
#   iwr https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.ps1 | iex
#   or
#   .\install.ps1 -Version v1.0.0 -Prefix "C:\Program Files\cbwsh"
#
# Options:
#   -Version <version>  Install a specific version (default: latest)
#   -Prefix <path>      Install to a custom location (default: $env:LOCALAPPDATA\cbwsh)
#   -AddToPath          Add installation directory to PATH
#   -Help               Show this help message
#

param(
    [string]$Version = "latest",
    [string]$Prefix = "$env:LOCALAPPDATA\cbwsh",
    [switch]$AddToPath = $false,
    [switch]$Help = $false
)

$ErrorActionPreference = "Stop"

# Configuration
$GithubRepo = "cbwinslow/cbwsh"
$BinaryName = "cbwsh"

function Show-Banner {
    Write-Host @"

   _____ ______          _____  _    _ 
  / ____|  _ \ \        / / __|| |  | |
 | |    | |_) \ \  /\  / / |__ | |__| |
 | |    |  _ < \ \/  \/ /\__ \|  __  |
 | |____| |_) | \  /\  / ___) | |  | |
  \_____|____/   \/  \/ |____/|_|  |_|
                                        
  Custom Bubble Tea Shell

"@ -ForegroundColor Cyan
}

function Show-Help {
    Write-Host @"
cbwsh PowerShell Installer

Usage: .\install.ps1 [options]

Options:
  -Version <version>  Install a specific version (default: latest)
  -Prefix <path>      Install to a custom location (default: `$env:LOCALAPPDATA\cbwsh)
  -AddToPath          Add installation directory to PATH
  -Help               Show this help message

Examples:
  .\install.ps1                           # Install latest version
  .\install.ps1 -Version v1.0.0           # Install specific version
  .\install.ps1 -Prefix "C:\Tools\cbwsh"  # Install to custom location
  .\install.ps1 -AddToPath                # Install and add to PATH

"@
}

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] " -ForegroundColor Blue -NoNewline
    Write-Host $Message
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARN] " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] " -ForegroundColor Red -NoNewline
    Write-Host $Message
    exit 1
}

function Get-Architecture {
    $arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")
    switch ($arch) {
        "AMD64" { return "amd64" }
        "x86" { return "386" }
        "ARM64" { return "arm64" }
        default { Write-Error "Unsupported architecture: $arch" }
    }
}

function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$GithubRepo/releases/latest" -Method Get
        $version = $response.tag_name
        
        if ([string]::IsNullOrEmpty($version)) {
            Write-Error "Could not determine latest version"
        }
        
        Write-Info "Latest version: $version"
        return $version
    }
    catch {
        Write-Error "Failed to fetch latest version: $_"
    }
}

function Download-Binary {
    param(
        [string]$Version,
        [string]$Architecture
    )
    
    $fileName = "${BinaryName}_windows_${Architecture}.zip"
    $downloadUrl = "https://github.com/$GithubRepo/releases/download/$Version/$fileName"
    $tempDir = [System.IO.Path]::GetTempPath()
    $tempFile = Join-Path $tempDir $fileName
    $extractDir = Join-Path $tempDir "cbwsh-extract"
    
    Write-Info "Downloading cbwsh $Version for windows/$Architecture..."
    Write-Info "URL: $downloadUrl"
    
    try {
        # Download the file
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
        
        # Create extraction directory
        if (Test-Path $extractDir) {
            Remove-Item -Recurse -Force $extractDir
        }
        New-Item -ItemType Directory -Path $extractDir -Force | Out-Null
        
        # Extract the archive
        Write-Info "Extracting archive..."
        Expand-Archive -Path $tempFile -DestinationPath $extractDir -Force
        
        # Find the binary
        $binaryPath = Get-ChildItem -Path $extractDir -Recurse -Filter "${BinaryName}.exe" | Select-Object -First 1
        
        if ($null -eq $binaryPath) {
            Write-Error "Binary not found in archive"
        }
        
        Write-Success "Download complete"
        return $binaryPath.FullName
    }
    catch {
        Write-Error "Download failed: $_"
    }
    finally {
        # Cleanup temp file
        if (Test-Path $tempFile) {
            Remove-Item $tempFile -Force -ErrorAction SilentlyContinue
        }
    }
}

function Install-Binary {
    param(
        [string]$BinaryPath,
        [string]$Prefix
    )
    
    Write-Info "Installing to $Prefix..."
    
    # Create prefix directory if it doesn't exist
    if (-not (Test-Path $Prefix)) {
        New-Item -ItemType Directory -Path $Prefix -Force | Out-Null
    }
    
    $destination = Join-Path $Prefix "${BinaryName}.exe"
    
    try {
        Copy-Item -Path $BinaryPath -Destination $destination -Force
        Write-Success "cbwsh installed to $destination"
        return $destination
    }
    catch {
        Write-Error "Failed to copy binary: $_"
    }
}

function Add-ToPath {
    param([string]$Prefix)
    
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    if ($currentPath -notlike "*$Prefix*") {
        Write-Info "Adding $Prefix to PATH..."
        
        $newPath = "$currentPath;$Prefix"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        
        # Update current session
        $env:Path = "$env:Path;$Prefix"
        
        Write-Success "Added to PATH"
    }
    else {
        Write-Info "Directory already in PATH"
    }
}

function New-DefaultConfig {
    $configDir = Join-Path $env:USERPROFILE ".cbwsh"
    $configFile = Join-Path $configDir "config.yaml"
    
    if (-not (Test-Path $configDir)) {
        Write-Info "Creating configuration directory..."
        New-Item -ItemType Directory -Path $configDir -Force | Out-Null
    }
    
    if (-not (Test-Path $configFile)) {
        Write-Info "Creating default configuration..."
        
        $configContent = @"
# cbwsh Configuration
# https://github.com/cbwinslow/cbwsh

shell:
  default_shell: powershell
  history_size: 10000

ui:
  theme: default
  layout: single
  show_status_bar: true
  enable_animations: true
  syntax_highlighting: true

ai:
  provider: none  # Options: none, openai, anthropic, gemini, local
  api_key: ""
  model: ""
  enable_suggestions: false

secrets:
  store_path: ~/.cbwsh/secrets.enc
  encryption_algorithm: AES-256-GCM
  key_derivation: argon2id

keybindings:
  quit: ctrl+q
  help: ctrl+?
  ai_assist: ctrl+a
"@
        
        $configContent | Out-File -FilePath $configFile -Encoding utf8
        Write-Success "Created default configuration at $configFile"
    }
    else {
        Write-Info "Configuration file already exists at $configFile"
    }
}

function Test-Installation {
    param([string]$Prefix)
    
    Write-Info "Verifying installation..."
    
    $binaryPath = Join-Path $Prefix "${BinaryName}.exe"
    
    if (Test-Path $binaryPath) {
        try {
            $version = & $binaryPath --version 2>&1
            Write-Success "cbwsh is installed and ready to use!"
        }
        catch {
            Write-Success "cbwsh is installed!"
        }
    }
    else {
        Write-Warning "Binary not found at $binaryPath"
    }
}

function Show-Instructions {
    Write-Host ""
    Write-Host "Installation Complete!" -ForegroundColor Green -BackgroundColor Black
    Write-Host ""
    Write-Host "To start cbwsh, run:"
    Write-Host ""
    Write-Host "  cbwsh" -ForegroundColor White
    Write-Host ""
    Write-Host "Quick Start:"
    Write-Host "  - Press Ctrl+? or F1 for help"
    Write-Host "  - Press Ctrl+A to toggle AI assist mode"
    Write-Host "  - Press Ctrl+Q to quit"
    Write-Host ""
    Write-Host "Configuration: $env:USERPROFILE\.cbwsh\config.yaml"
    Write-Host ""
    Write-Host "For more information, visit:"
    Write-Host "  https://github.com/$GithubRepo"
    Write-Host ""
    
    if (-not $AddToPath) {
        Write-Host "Note: To add cbwsh to your PATH, run:" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "  .\install.ps1 -AddToPath" -ForegroundColor White
        Write-Host ""
    }
}

function Main {
    Show-Banner
    
    if ($Help) {
        Show-Help
        return
    }
    
    # Get architecture
    $arch = Get-Architecture
    Write-Info "Detected platform: windows/$arch"
    
    # Get version
    if ($Version -eq "latest") {
        $Version = Get-LatestVersion
    }
    
    # Download binary
    $binaryPath = Download-Binary -Version $Version -Architecture $arch
    
    # Install binary
    $installedPath = Install-Binary -BinaryPath $binaryPath -Prefix $Prefix
    
    # Add to PATH if requested
    if ($AddToPath) {
        Add-ToPath -Prefix $Prefix
    }
    
    # Create default config
    New-DefaultConfig
    
    # Verify installation
    Test-Installation -Prefix $Prefix
    
    # Show instructions
    Show-Instructions
    
    # Cleanup
    $extractDir = Join-Path ([System.IO.Path]::GetTempPath()) "cbwsh-extract"
    if (Test-Path $extractDir) {
        Remove-Item -Recurse -Force $extractDir -ErrorAction SilentlyContinue
    }
}

# Run main function
Main
