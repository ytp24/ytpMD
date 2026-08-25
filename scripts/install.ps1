$ErrorActionPreference = "Stop"

$Version = "3.2.0"
$InstallDir = Join-Path $env:LOCALAPPDATA "ytpMD\bin"

Write-Host "[*] Installing ytpMD (v$Version) for Windows..." -ForegroundColor Cyan

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

$BinaryPath = Join-Path $InstallDir "ytpMD.exe"
$Alias1 = Join-Path $InstallDir "ytp24.exe"
$Alias2 = Join-Path $InstallDir "pdf2md.exe"

if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "[*] Compiling binary with Go..." -ForegroundColor Cyan
    go build -ldflags="-s -w" -o $BinaryPath ./cmd/ytpMD
    Copy-Item -Path $BinaryPath -Destination $Alias1 -Force
    Copy-Item -Path $BinaryPath -Destination $Alias2 -Force
} else {
    Write-Host "[x] Go is required for local build. Please install Go from https://golang.org/" -ForegroundColor Red
    Exit 1
}

$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "[*] Adding $InstallDir to User PATH..." -ForegroundColor Cyan
    [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    $env:PATH = "$env:PATH;$InstallDir"
}

Write-Host "[+] Successfully installed ytpMD to $BinaryPath" -ForegroundColor Green
Write-Host "You can now open a new PowerShell terminal and type 'ytpMD' to launch the tool." -ForegroundColor Green
