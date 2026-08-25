# ==============================================================================
# ytpMD [pdf2md] - Official Windows PowerShell Installer
# ==============================================================================
$ErrorActionPreference = "Stop"

$Version = "3.1.0"
$InstallDir = Join-Path $env:LOCALAPPDATA "ytp24\bin"

Write-Host "[*] Installing ytpMD (ytp24 v$Version) for Windows..." -ForegroundColor Cyan

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

$BinaryPath = Join-Path $InstallDir "ytp24.exe"
$AliasPath = Join-Path $InstallDir "pdf2md.exe"

# Build if Go is available
if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "[*] Compiling binary with Go..." -ForegroundColor Cyan
    go build -ldflags="-s -w" -o $BinaryPath ./cmd/ytp24
    Copy-Item -Path $BinaryPath -Destination $AliasPath -Force
} else {
    Write-Host "[x] Go is required for local build. Please install Go from https://golang.org/" -ForegroundColor Red
    Exit 1
}

# Add to User PATH if not already present
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "[*] Adding $InstallDir to User PATH..." -ForegroundColor Cyan
    [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    $env:PATH = "$env:PATH;$InstallDir"
}

Write-Host "[+] Successfully installed ytp24 to $BinaryPath" -ForegroundColor Green
Write-Host "You can now open a new PowerShell terminal and type 'ytp24' to launch the tool." -ForegroundColor Green
