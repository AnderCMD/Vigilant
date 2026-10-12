<#
.SYNOPSIS
  Builds Vigilant into dist/ (Windows exe with icon and version info).
.PARAMETER Version
  Semantic version to embed, e.g. 1.2.3. Defaults to "dev".
.PARAMETER Arch
  Target architecture: amd64 (default) or arm64.
#>
param(
  [string]$Version = "dev",
  [ValidateSet("amd64", "arm64")][string]$Arch = "amd64"
)
$ErrorActionPreference = "Stop"
Set-Location (Split-Path $PSScriptRoot -Parent)

$numeric = if ($Version -match '^\d+\.\d+\.\d+') { $Matches[0] } else { "0.0.0" }

if (-not (Get-Command go-winres -ErrorAction SilentlyContinue)) {
  Write-Host "Installing go-winres..."
  go install github.com/tc-hib/go-winres@v0.3.3
  $env:PATH += ";" + (go env GOPATH) + "\bin"
}

Push-Location cmd/vigilant
go-winres make --arch $Arch --product-version "$numeric.0" --file-version "$numeric.0"
Pop-Location

$env:GOOS = "windows"; $env:GOARCH = $Arch; $env:CGO_ENABLED = "0"
New-Item -ItemType Directory -Force dist | Out-Null
$out = "dist/vigilant-windows-$Arch.exe"
go build -trimpath -ldflags "-s -w -H=windowsgui -X main.version=$Version" -o $out ./cmd/vigilant
Write-Host "Built $out"
