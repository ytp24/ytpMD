$ErrorActionPreference = "Stop"

$Version = "3.2.0"
$InstallDir = Join-Path $env:LOCALAPPDATA "ytpMD\bin"
$Repo = "ytp24/ytpMD"

Write-Host "====================================================" -ForegroundColor Cyan
Write-Host "   ytpMD [pdf2md] Installer for Windows (v$Version)" -ForegroundColor Cyan
Write-Host "====================================================" -ForegroundColor Cyan
Write-Host ""

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

$BinaryPath = Join-Path $InstallDir "ytpmd.exe"
$McpPath = Join-Path $InstallDir "ytpmd-mcp.exe"
$Alias1 = Join-Path $InstallDir "ytpMD.exe"
$Alias2 = Join-Path $InstallDir "ytp24.exe"
$Alias3 = Join-Path $InstallDir "pdf2md.exe"

$Installed = $false

if (Test-Path ".\go.mod" -and (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "[*] Compiling binaries from local repository source with Go..." -ForegroundColor Cyan
    go build -ldflags="-s -w" -o $BinaryPath ./cmd/ytpmd
    go build -ldflags="-s -w" -o $McpPath ./cmd/ytpmd-mcp
    $Installed = $true
} else {
    $ZipName = "ytpMD-$Version-windows-amd64.zip"
    $DownloadUrl = "https://github.com/$Repo/releases/download/v$Version/$ZipName"
    $TempZip = Join-Path $env:TEMP $ZipName

    Write-Host "[*] Downloading release archive for Windows ($ZipName)..." -ForegroundColor Cyan
    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempZip -UseBasicParsing
        if (Test-Path $TempZip) {
            $ExtractDir = Join-Path $env:TEMP "ytpmd_extract_$PID"
            Expand-Archive -Path $TempZip -DestinationPath $ExtractDir -Force
            
            $ExtractedBin = Join-Path $ExtractDir "ytpMD-$Version-windows-amd64\ytpmd.exe"
            if (-not (Test-Path $ExtractedBin)) {
                $ExtractedBin = Join-Path $ExtractDir "ytpMD-$Version-windows-amd64\ytpMD.exe"
            }
            Copy-Item -Path $ExtractedBin -Destination $BinaryPath -Force

            $ExtractedMcp = Join-Path $ExtractDir "ytpMD-$Version-windows-amd64\ytpmd-mcp.exe"
            if (Test-Path $ExtractedMcp) {
                Copy-Item -Path $ExtractedMcp -Destination $McpPath -Force
            }

            Remove-Item -Path $ExtractDir -Recurse -Force -ErrorAction SilentlyContinue
            Remove-Item -Path $TempZip -Force -ErrorAction SilentlyContinue
            $Installed = $true
        }
    } catch {
        Write-Host "[!] Prebuilt zip download failed. Trying 'go install'..." -ForegroundColor Yellow
        if (Get-Command go -ErrorAction SilentlyContinue) {
            go install "github.com/$Repo/cmd/ytpmd@latest"
            $GopathBin = Join-Path (go env GOPATH) "bin\ytpmd.exe"
            if (Test-Path $GopathBin) {
                Copy-Item -Path $GopathBin -Destination $BinaryPath -Force
                $Installed = $true
            }
        }
    }
}

if ($Installed) {
    Copy-Item -Path $BinaryPath -Destination $Alias1 -Force
    Copy-Item -Path $BinaryPath -Destination $Alias2 -Force
    Copy-Item -Path $BinaryPath -Destination $Alias3 -Force

    $UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        Write-Host "[*] Adding $InstallDir to User PATH..." -ForegroundColor Cyan
        [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
        $env:PATH = "$env:PATH;$InstallDir"
    }

    Write-Host ""
    Write-Host "[+] Successfully installed ytpMD to $InstallDir" -ForegroundColor Green
    Write-Host "    Available commands: ytpmd, ytpMD, ytp24, pdf2md, ytpmd-mcp" -ForegroundColor Green
    Write-Host ""
    Write-Host "Open a new PowerShell window and run 'ytpmd' to start!" -ForegroundColor Cyan
} else {
    Write-Host "[x] Installation could not complete. Please install Go or download manually." -ForegroundColor Red
    Exit 1
}
