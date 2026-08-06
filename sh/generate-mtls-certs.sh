#!/bin/bash

# generate-mtls-certs.sh

# Параметры по умолчанию
CertsPath="${1:-../secrets/certs}"
DaysValid="${2:-365}"

# Функция для генерации сертификата сервиса
new_service_certificate() {
    local ServiceName="$1"
    local CertsRootPath="$2"
    local CAKeyPath="$3"
    local CACertPath="$4"
    local DaysValid="$5"
    
    local ServicePath="$CertsRootPath/$ServiceName"
    if [[ ! -d "$ServicePath" ]]; then
        mkdir -p "$ServicePath"
    fi
    
    local ServiceKeyPath="$ServicePath/server-key.pem"
    openssl genrsa -out "$ServiceKeyPath" 4096 2>/dev/null
    
    # Создаем временный конфиг с SAN (Subject Alternative Names)
    local SanConfigPath="$ServicePath/san.cnf"
    cat <<EOF > "$SanConfigPath"
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
EOF
    
    local CsrPath="$ServicePath/server.csr"
    openssl req -new -key "$ServiceKeyPath" -out "$CsrPath" -config "$SanConfigPath" 2>/dev/null
    
    local ServiceCertPath="$ServicePath/server-cert.pem"
    if openssl x509 -req -in "$CsrPath" -CA "$CACertPath" -CAkey "$CAKeyPath" \
        -CAcreateserial -out "$ServiceCertPath" -days "$DaysValid" -sha256 \
        -extfile "$SanConfigPath" -extensions req_ext 2>/dev/null; then
        
        rm -f "$CsrPath" "$SanConfigPath"
        
        cp -f "$CACertPath" "$ServicePath/ca-cert.pem"
        return 0
    fi
    return 1
}

# Функция для генерации клиентского сертификата
new_client_certificate() {
    local ClientName="$1"
    local ClientPath="$2"
    local CAKeyPath="$3"
    local CACertPath="$4"
    local DaysValid="$5"
    
    local ClientKeyPath="$ClientPath/client-key.pem"
    openssl genrsa -out "$ClientKeyPath" 4096 2>/dev/null
    
    local ClientCsrPath="$ClientPath/client.csr"
    openssl req -new -key "$ClientKeyPath" -out "$ClientCsrPath" -subj "/CN=$ClientName" 2>/dev/null
    
    local ClientCertPath="$ClientPath/client-cert.pem"
    if openssl x509 -req -in "$ClientCsrPath" -CA "$CACertPath" -CAkey "$CAKeyPath" \
        -CAcreateserial -out "$ClientCertPath" -days "$DaysValid" -sha256 2>/dev/null; then
        
        rm -f "$ClientCsrPath"
        return 0
    fi
    return 1
}

# --- Основной скрипт ---

# Проверка наличия openssl
if ! command -v openssl &> /dev/null; then
    echo -e "\033[31mERROR: OpenSSL not found\033[0m"
    exit 1
fi

# Создаем корневую папку для сертификатов
if [[ ! -d "$CertsPath" ]]; then
    mkdir -p "$CertsPath"
fi

CAKeyPath="$CertsPath/ca-key.pem"
CACertPath="$CertsPath/ca-cert.pem"

# Генерация корневого CA, если отсутствует
if [[ ! -f "$CAKeyPath" ]] || [[ ! -f "$CACertPath" ]]; then
    openssl genrsa -out "$CAKeyPath" 4096 2>/dev/null
    openssl req -x509 -new -nodes -key "$CAKeyPath" -sha256 -days 3650 \
        -out "$CACertPath" -subj "/CN=Internal CA for Microservices" 2>/dev/null
fi

echo -e "\033[33mWARNING: ca-key.pem must be stored securely OFFLINE! Never commit to git.\033[0m"

# Массив сервисов
services=("article-service" "user-service" "simulator-service")

for service in "${services[@]}"; do
    new_service_certificate "$service" "$CertsPath" "$CAKeyPath" "$CACertPath" "$DaysValid"
done

# Генерация сертификатов для Gateway (клиентские для mTLS)
GatewayPath="$CertsPath/gateway-server"
if [[ ! -d "$GatewayPath" ]]; then
    mkdir -p "$GatewayPath"
fi

new_client_certificate "gateway-client" "$GatewayPath" "$CAKeyPath" "$CACertPath" "$DaysValid"

# Генерация публичных сертификатов для Gateway (внешний HTTPS)
GatewayPublicKeyPath="$GatewayPath/public-key.pem"
GatewayPublicCertPath="$GatewayPath/public-cert.pem"

openssl genrsa -out "$GatewayPublicKeyPath" 4096 2>/dev/null
openssl req -new -key "$GatewayPublicKeyPath" -out "$GatewayPath/public.csr" -subj "/CN=gateway-public" 2>/dev/null
openssl x509 -req -in "$GatewayPath/public.csr" -CA "$CACertPath" -CAkey "$CAKeyPath" \
    -CAcreateserial -out "$GatewayPublicCertPath" -days "$DaysValid" -sha256 2>/dev/null
rm -f "$GatewayPath/public.csr"

# Копирование CA сертификата в папку gateway
cp -f "$CACertPath" "$GatewayPath/ca-cert.pem"

echo -e "\033[32mGateway certs generated\033[0m"
echo -e "\033[32marticle-service certs generated\033[0m"
echo -e "\033[32muser-service certs generated\033[0m"
echo -e "\033[32msimulator-service certs generated\033[0m"