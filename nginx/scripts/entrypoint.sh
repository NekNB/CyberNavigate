#!/bin/bash
set -e

# Пути к сертификатам
CERT_PATH="/etc/letsencrypt/live/xn----8sbabhgkwfn1brsi2a.xn--p1ai/fullchain.pem"

mkdir -p /var/www/certbot

# Проверяем, есть ли уже сертификаты ДО запуска Nginx
if [ ! -f "$CERT_PATH" ]; then
    echo "Сертификаты не найдены. Используем HTTP конфиг..."
    cp /etc/nginx/configs/nginx-http.conf /etc/nginx/conf.d/default.conf
else
    echo "Сертификаты существуют. Используем HTTPS конфиг..."
    cp /etc/nginx/configs/nginx-https.conf /etc/nginx/conf.d/default.conf
fi

# Запускаем Nginx в фоновом режиме
echo "Запуск Nginx..."
nginx -g 'daemon off;' &

# Ждем, пока Nginx поднимется
sleep 3
# Если сертификатов нет, получаем их
if [ ! -f "$CERT_PATH" ]; then
    echo "Запрашиваем сертификаты у Let's Encrypt..."
    certbot certonly --webroot -w /var/www/certbot \
        --email admin@xn----8sbabhgkwfn1brsi2a.xn--p1ai \
        --agree-tos \
        --no-eff-email \
        --non-interactive \
        -d xn----8sbabhgkwfn1brsi2a.xn--p1ai \
        -d api.xn----8sbabhgkwfn1brsi2a.xn--p1ai \
        -d www.xn----8sbabhgkwfn1brsi2a.xn--p1ai

    # Заменяем конфиг Nginx на версию с HTTPS
    echo "Активируем HTTPS конфигурацию..."
    cp /etc/nginx/configs/nginx-https.conf /etc/nginx/conf.d/default.conf

    # Перезагружаем Nginx
    echo "Перезагрузка Nginx..."
    nginx -s reload
fi

# Цикл автообновления
while true; do
    sleep 12h
    echo "Проверка необходимости обновления сертификатов..."
    certbot renew --quiet
    nginx -s reload
done
