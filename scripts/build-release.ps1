$ErrorActionPreference = "Stop"

$projectDirectory = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$releaseDirectory = Join-Path $projectDirectory "release"
$executablePath = Join-Path $releaseDirectory "TestDataGenerator.exe"
$archivePath = Join-Path $releaseDirectory "test-data-generator-windows.zip"

New-Item -ItemType Directory -Force -Path $releaseDirectory | Out-Null

Push-Location $projectDirectory
try {
    go test ./...
    if ($LASTEXITCODE -ne 0) {
        throw "Tests failed; release archive was not updated."
    }

    go build -o $executablePath .\cmd\test-data-generator
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed; release archive was not updated."
    }

    Compress-Archive -LiteralPath $executablePath -DestinationPath $archivePath -Force
} finally {
    Pop-Location
}

Write-Host "Ready: $archivePath"
