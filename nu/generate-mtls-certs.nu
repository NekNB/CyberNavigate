#!/bin/bash

# generate-mtls-certs.sh



# --- Основной скрипт ---
def main [
    CertsPath?: string = "../secrets/certs"
    DaysValid?: string = "365"
] {

    # Проверка наличия openssl
    if (which openssl | is-empty) {
        print $"(ansi r)ERROR: OpenSSL not found"
        exit 1
    }

    # Создаем корневую папку для сертификатов
    mkdir $CertsPath

    let CAKeyPath = $"($CertsPath)/ca-key.pem"
    let CACertPath = $"($CertsPath)/ca-cert.pem"

    # Генерация корневого CA, если отсутствует
    if not ($CAKeyPath | path exists) or not ($CACertPath | path exists) {
        ^openssl genrsa -out $CAKeyPath 4096 | ignore --stderr 
        ^openssl req -x509 -new -nodes -key $CAKeyPath -sha256 -days 3650 -out $CACertPath -subj "/CN=Internal CA for Microservices" | ignore --stderr 
    }

    print $"(ansi y)WARNING: ca-key.pem must be stored securely OFFLINE! Never commit to git"

    # Массив сервисов
    let services = ["article-service" "user-service" "simulator-service"]

    for service in $services {
        new-service-certificate $service $CertsPath $CAKeyPath $CACertPath $DaysValid
    }
    

    # Генерация сертификатов для Gateway (клиентские для mTLS)
    let GatewayPath = $"($CertsPath)/gateway-server"
    mkdir $GatewayPath
    
    new-client-certificate "gateway-client" $GatewayPath $CAKeyPath $CACertPath $DaysValid

    # Генерация публичных сертификатов для Gateway (внешний HTTPS)
    let GatewayPublicKeyPath = $"($GatewayPath)/public-key.pem"
    let GatewayPublicCertPath = $"($GatewayPath)/public-cert.pem"

    ^openssl genrsa -out $GatewayPublicKeyPath 4096 | ignore --stderr 
    ^openssl req -new -key $GatewayPublicKeyPath -out $"($GatewayPath)/public.csr" -subj "/CN=gateway-public" | ignore --stderr 
    ^openssl x509 -req -in $"($GatewayPath)/public.csr" -CA $CACertPath -CAkey $CAKeyPath -CAcreateserial -out $GatewayPublicCertPath -days $DaysValid -sha256 | ignore --stderr 
    rm -f $"($GatewayPath)/public.csr" $"($GatewayPath)/client.csr"

    # Копирование CA сертификата в папку gateway
    cp -f $CACertPath $"($GatewayPath)/ca-cert.pem"
    rm -f $"($CertsPath)/ca-cert.srl"

    print $"(ansi blue)Gateway certs (ansi green)generated"
    print $"(ansi blue)Article-service certs (ansi green)generated"
    print $"(ansi blue)User-service certs (ansi green)generated"
    print $"(ansi blue)Simulator-service certs (ansi green)generated"
}


def new-service-certificate [
    ServiceName: string
    CertsRootPath: string
    CAKeyPath: string
    CACertPath: string
    DaysValid: string
] {
    let ServicePath = $"($CertsRootPath)/($ServiceName)"
    mkdir $ServicePath
   

    let ServiceKeyPath = $"($ServicePath)/server-key.pem"
    openssl genrsa -out $ServiceKeyPath 4096 | ignore --stderr 
    
    # Создаем временный конфиг с SAN (Subject Alternative Names)
    let SanConfigPath = $"($ServicePath)/san.cnf"

$"[req]
default_bits = 4096
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = req_ext

[dn]
CN=($ServiceName)

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = ($ServiceName)
DNS.2 = localhost
DNS.3 = *.local
    " | save -f $SanConfigPath
    let CsrPath = $"($ServicePath)/server.csr"
    ^openssl req -new -key $ServiceKeyPath -out $CsrPath -config $SanConfigPath | ignore --stderr 

    let ServiceCertPath = $"($ServicePath)/server-cert.pem"
    try { 
        ^openssl x509 -req -in $CsrPath -CA $CACertPath -CAkey $CAKeyPath -CAcreateserial -out $ServiceCertPath -days $DaysValid -sha256 -extfile $SanConfigPath -extensions req_ext | ignore --stderr 
        return 1
    } catch {|err|
        rm -f $CsrPath $SanConfigPath
        cp -f $CACertPath $"($ServicePath)/ca-cert.pem"
        return 0
    }

}

# Функция для генерации клиентского сертификата
def new-client-certificate [
    ClientName: string
    ClientPath: string
    CAKeyPath: string
    CACertPath: string
    DaysValid: string
] {    
    let ClientKeyPath = $"($ClientPath)/client-key.pem"
    ^openssl genrsa -out $ClientKeyPath 4096 | ignore --stderr 
    
    let ClientCsrPath = $"($ClientPath)/client.csr"
    ^openssl req -new -key $ClientKeyPath -out $ClientCsrPath -subj $"/CN=($ClientName)" | ignore --stderr 
    
    let ClientCertPath = $"($ClientPath)/client-cert.pem"
    try {
        ^openssl x509 -req -in $ClientCsrPath -CA $CACertPath -CAkey $CAKeyPath -CAcreateserial -out $ClientCertPath -days $DaysValid -sha256 | ignore --stderr 
        return 1
    } catch {
        |err|
        rm -f $ClientCsrPath
        return 0
    }
}

