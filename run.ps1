$ErrorActionPreference = 'Stop'

Write-Host 'Building application before running...'
.\build.ps1

Write-Host 'Starting application...'
.\booking-system.exe