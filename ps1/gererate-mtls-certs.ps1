# generate-mtls-certs.ps1

param(
    [string]$CertsPath = "../secrets/certs",
    [int]$DaysValid = 365
)

function New-ServiceCertificate {
    param(
        [string]$ServiceName,
        [string]$CertsRootPath,
        [string]$CAKeyPath,
        [string]$CACertPath,
        [int]$DaysValid = 365
    )
    
    $ServicePath = Join-Path $CertsRootPath $ServiceName
    if (-not (Test-Path $ServicePath)) {
        New-Item -ItemType Directory -Path $ServicePath -Force | Out-Null
    }
    
    $ServiceKeyPath = Join-Path $ServicePath "server-key.pem"
    openssl genrsa -out $ServiceKeyPath 4096 2>$null
    
    $SanConfig = @"
[req]
default_bits = 4096
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = req_ext

[dn]
CN=$ServiceName

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = $ServiceName
DNS.2 = localhost
DNS.3 = *.local
"@
    
    $SanConfigPath = Join-Path $ServicePath "san.cnf"
    $SanConfig | Out-File -FilePath $SanConfigPath -Encoding ASCII
    
    $CsrPath = Join-Path $ServicePath "server.csr"
    openssl req -new -key $ServiceKeyPath -out $CsrPath -config $SanConfigPath 2>$null
    
    $ServiceCertPath = Join-Path $ServicePath "server-cert.pem"
    openssl x509 -req -in $CsrPath -CA $CACertPath -CAkey $CAKeyPath `
        -CAcreateserial -out $ServiceCertPath -days $DaysValid -sha256 `
        -extfile $SanConfigPath -extensions req_ext 2>$null
    
    if ($LASTEXITCODE -eq 0) {
        Remove-Item $CsrPath -Force -ErrorAction SilentlyContinue
        Remove-Item $SanConfigPath -Force -ErrorAction SilentlyContinue
        
        $ServiceCACopy = Join-Path $ServicePath "ca-cert.pem"
        Copy-Item $CACertPath -Destination $ServiceCACopy -Force
        return $true
    }
    return $false
}

function New-ClientCertificate {
    param(
        [string]$ClientName,
        [string]$ClientPath,
        [string]$CAKeyPath,
        [string]$CACertPath,
        [int]$DaysValid = 365
    )
    
    $ClientKeyPath = Join-Path $ClientPath "client-key.pem"
    openssl genrsa -out $ClientKeyPath 4096 2>$null
    
    $ClientCsrPath = Join-Path $ClientPath "client.csr"
    $ClientSubject = "/CN=$ClientName"
    openssl req -new -key $ClientKeyPath -out $ClientCsrPath -subj $ClientSubject 2>$null
    
    $ClientCertPath = Join-Path $ClientPath "client-cert.pem"
    openssl x509 -req -in $ClientCsrPath -CA $CACertPath -CAkey $CAKeyPath `
        -CAcreateserial -out $ClientCertPath -days $DaysValid -sha256 2>$null
    
    if ($LASTEXITCODE -eq 0) {
        Remove-Item $ClientCsrPath -Force -ErrorAction SilentlyContinue
        return $true
    }
    return $false
}

# Main script
$opensslCheck = Get-Command openssl -ErrorAction SilentlyContinue
if (-not $opensslCheck) {
    Write-Host "ERROR: OpenSSL not found" -ForegroundColor Red
    exit 1
}

if (-not (Test-Path $CertsPath)) {
    New-Item -ItemType Directory -Path $CertsPath -Force | Out-Null
}

$CAKeyPath = Join-Path $CertsPath "ca-key.pem"
$CACertPath = Join-Path $CertsPath "ca-cert.pem"

if (-not (Test-Path $CAKeyPath) -or -not (Test-Path $CACertPath)) {
    openssl genrsa -out $CAKeyPath 4096 2>$null
    openssl req -x509 -new -nodes -key $CAKeyPath -sha256 -days 3650 `
        -out $CACertPath -subj "/CN=Internal CA for Microservices" 2>$null
}

Write-Host "WARNING: ca-key.pem must be stored securely OFFLINE! Never commit to git." -ForegroundColor Yellow

$services = @("article-service", "user-service", "simulator-service")

foreach ($service in $services) {
    New-ServiceCertificate -ServiceName $service -CertsRootPath $CertsPath -CAKeyPath $CAKeyPath -CACertPath $CACertPath -DaysValid $DaysValid
}

# Generate certificates for Gateway (client certificates for mTLS)
$GatewayPath = Join-Path $CertsPath "gateway-server"
if (-not (Test-Path $GatewayPath)) {
    New-Item -ItemType Directory -Path $GatewayPath -Force | Out-Null
}

New-ClientCertificate -ClientName "gateway-client" -ClientPath $GatewayPath -CAKeyPath $CAKeyPath -CACertPath $CACertPath -DaysValid $DaysValid

# Generate public certificates for Gateway (external HTTPS)
$GatewayPublicKeyPath = Join-Path $GatewayPath "public-key.pem"
$GatewayPublicCertPath = Join-Path $GatewayPath "public-cert.pem"

openssl genrsa -out $GatewayPublicKeyPath 4096 2>$null
openssl req -new -key $GatewayPublicKeyPath -out "$GatewayPath\public.csr" -subj "/CN=gateway-public" 2>$null
openssl x509 -req -in "$GatewayPath\public.csr" -CA $CACertPath -CAkey $CAKeyPath `
    -CAcreateserial -out $GatewayPublicCertPath -days $DaysValid -sha256 2>$null
Remove-Item "$GatewayPath\public.csr" -Force -ErrorAction SilentlyContinue

# Copy CA cert to gateway folder
Copy-Item $CACertPath -Destination (Join-Path $GatewayPath "ca-cert.pem") -Force

Write-Host "Gateway certs generated" -ForegroundColor Green
Write-Host "article-service certs generated" -ForegroundColor Green
Write-Host "user-service certs generated" -ForegroundColor Green