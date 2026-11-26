$ErrorActionPreference = 'Stop'
$env:CGO_ENABLED = '0'
Write-Host "Compilando BrightVault.exe..." -ForegroundColor Cyan
go build -ldflags="-H=windowsgui" -o BrightVault.exe ./cmd/ytcloud
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
Write-Host "BrightVault.exe compilado com sucesso!" -ForegroundColor Green
