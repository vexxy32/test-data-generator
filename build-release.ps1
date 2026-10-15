$ErrorActionPreference = "Stop"

$projectDirectory = Split-Path -Parent $MyInvocation.MyCommand.Path
$releaseDirectory = Join-Path $projectDirectory "release"
$executablePath = Join-Path $releaseDirectory "TestDataGenerator.exe"
$archivePath = Join-Path $releaseDirectory "test-data-generator-windows.zip"

New-Item -ItemType Directory -Force -Path $releaseDirectory | Out-Null

Push-Location $projectDirectory
try {
    go build -o $executablePath .
    Compress-Archive -LiteralPath $executablePath -DestinationPath $archivePath -Force
} finally {
    Pop-Location
}

Write-Host "Ready: $archivePath"
