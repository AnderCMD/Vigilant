# Runs Vigilant from source (no console window, so output is hidden).
Set-Location (Split-Path $PSScriptRoot -Parent)
go run -ldflags "-H=windowsgui" ./cmd/vigilant
