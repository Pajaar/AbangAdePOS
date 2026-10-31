<#
.SYNOPSIS
  Menjalankan golang-migrate ke database yang ada di .env
  (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME).

.DESCRIPTION
  - .env dibaca dari root proyek, jadi skrip bisa dijalankan dari folder mana pun.
  - Password di-encode untuk URL dan disamarkan (***) di output.
  - Perintah yang bisa menurunkan/menghapus skema (drop, down, force, goto)
    minta konfirmasi dengan mengetik nama database.

.EXAMPLE
  .\scripts\migrate.ps1 up
  .\scripts\migrate.ps1 version
  .\scripts\migrate.ps1 down 1
#>
param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$MigrateArgs
)

$root = Split-Path -Parent $PSScriptRoot
$envFile = Join-Path $root '.env'

if (-not $MigrateArgs -or $MigrateArgs.Count -eq 0) {
    Write-Host 'Pemakaian: .\scripts\migrate.ps1 <perintah migrate>'
    Write-Host 'Contoh   : .\scripts\migrate.ps1 up | version | down 1'
    exit 1
}

if (-not (Test-Path $envFile)) {
    Write-Host ".env tidak ditemukan di $root. Salin dari .env.example lalu isi dulu." -ForegroundColor Red
    exit 1
}

if (-not (Get-Command migrate -ErrorAction SilentlyContinue)) {
    Write-Host 'Perintah migrate tidak ditemukan. Pastikan folder go\bin ada di PATH.' -ForegroundColor Red
    exit 1
}

# Baca .env tanpa mengubah environment sesi PowerShell.
$vars = @{}
Get-Content $envFile | ForEach-Object {
    if ($_ -match '^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$') {
        $vars[$matches[1]] = $matches[2].Trim().Trim('"')
    }
}
foreach ($key in 'DB_HOST', 'DB_PORT', 'DB_USER', 'DB_NAME') {
    if (-not $vars[$key]) {
        Write-Host "$key kosong di .env" -ForegroundColor Red
        exit 1
    }
}

$dbName = $vars['DB_NAME']
$target = "$dbName @ $($vars['DB_HOST']):$($vars['DB_PORT'])"

$perintah = $MigrateArgs[0].ToLower()
if ($perintah -in 'drop', 'down', 'force', 'goto') {
    Write-Host "PERINGATAN: 'migrate $($MigrateArgs -join ' ')' akan mengubah atau menghapus skema di $target" -ForegroundColor Yellow
    $jawab = Read-Host "Ketik nama database ($dbName) untuk lanjut"
    if ($jawab -ne $dbName) {
        Write-Host 'Dibatalkan.'
        exit 1
    }
}

$rawPass = [string]$vars['DB_PASSWORD']
$pass = [uri]::EscapeDataString($rawPass)
$url = "mysql://$($vars['DB_USER']):$pass@tcp($($vars['DB_HOST']):$($vars['DB_PORT']))/$dbName"

Write-Host "Database: $target"

# migrate di Windows tidak menerima path absolut (C:\...), jadi jalankan
# dari root proyek dengan path relatif.
# migrate menulis progres ke stderr; jangan sampai dianggap error PowerShell.
$ErrorActionPreference = 'Continue'
Push-Location $root
try {
    & migrate -path migrations -database $url @MigrateArgs 2>&1 | ForEach-Object {
        $line = "$_"
        if ($rawPass) {
            $line = $line.Replace($pass, '***').Replace($rawPass, '***')
        }
        $line
    }
    $code = $LASTEXITCODE
}
finally {
    Pop-Location
}
exit $code
