#Requires -Version 5.1
<#
.SYNOPSIS
    mmdc installer script for Windows
.DESCRIPTION
    Downloads and installs the latest version of mmdc
.EXAMPLE
    irm https://raw.githubusercontent.com/ageha734/mmdc/master/install.ps1 | iex
#>

$ErrorActionPreference = 'Stop'

$Repo = "ageha734/mmdc"
$InstallDir = "$env:LOCALAPPDATA\mmdc"

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Warn {
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

function Get-Platform {
    $arch = switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { "amd64" }
        "ARM64" { "arm64" }
        "x86"   { "386" }
        default { Write-Error "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
    }
    return "windows_$arch"
}

function Get-LatestVersion {
    $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
    return $response.tag_name
}

function Install-Mmdc {
    Write-Host ""
    Write-Host "mmdc installer"
    Write-Host "=============="
    Write-Host ""

    $platform = Get-Platform
    Write-Info "Detected platform: $platform"

    $version = Get-LatestVersion
    if (-not $version) {
        Write-Error "Failed to get latest version"
    }
    Write-Info "Latest version: $version"

    $filename = "mmdc_$platform.zip"
    $downloadUrl = "https://github.com/$Repo/releases/download/$version/$filename"
    $tempPath = Join-Path $env:TEMP $filename

    Write-Info "Downloading from: $downloadUrl"

    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempPath -UseBasicParsing
    }
    catch {
        Write-Error "Failed to download: $_"
    }

    $checksumsUrl = "https://github.com/$Repo/releases/download/$version/checksums.txt"
    $checksumsPath = Join-Path $env:TEMP "checksums.txt"

    try {
        Invoke-WebRequest -Uri $checksumsUrl -OutFile $checksumsPath -UseBasicParsing -ErrorAction SilentlyContinue
        Write-Info "Verifying checksum..."

        $checksums = Get-Content $checksumsPath
        $expectedHash = ($checksums | Where-Object { $_ -match $filename }) -split '\s+' | Select-Object -First 1
        $actualHash = (Get-FileHash -Path $tempPath -Algorithm SHA256).Hash.ToLower()

        if ($expectedHash -and $actualHash -ne $expectedHash) {
            Write-Error "Checksum verification failed"
        }
        Write-Info "Checksum verified"
        Remove-Item $checksumsPath -Force -ErrorAction SilentlyContinue
    }
    catch {
        Write-Warn "Could not verify checksum"
    }

    Write-Info "Extracting..."

    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $extractPath = Join-Path $env:TEMP "mmdc_extract"
    if (Test-Path $extractPath) {
        Remove-Item $extractPath -Recurse -Force
    }
    Expand-Archive -Path $tempPath -DestinationPath $extractPath -Force

    $binaryPath = Join-Path $extractPath "mmdc.exe"
    if (-not (Test-Path $binaryPath)) {
        $binaryPath = Get-ChildItem -Path $extractPath -Filter "mmdc.exe" -Recurse | Select-Object -First 1 -ExpandProperty FullName
    }

    Move-Item -Path $binaryPath -Destination (Join-Path $InstallDir "mmdc.exe") -Force

    Remove-Item $tempPath -Force -ErrorAction SilentlyContinue
    Remove-Item $extractPath -Recurse -Force -ErrorAction SilentlyContinue

    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$InstallDir*") {
        Write-Info "Adding $InstallDir to PATH..."
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallDir", "User")
        $env:Path = "$env:Path;$InstallDir"
    }

    Write-Info "Successfully installed mmdc to $InstallDir"
    Write-Host ""
    Write-Info "Run 'mmdc --version' to verify the installation"
    Write-Info "You may need to restart your terminal for PATH changes to take effect"
    Write-Host ""
    Write-Host "Installation complete!" -ForegroundColor Green
    Write-Host ""
}

Install-Mmdc
