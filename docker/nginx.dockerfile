FROM nginx:latest

# Устанавливаем certbot
RUN apt-get update && \
    apt-get install -y certbot python3-certbot-nginx openssl && \
    rm -rf /var/lib/apt/lists/*

# Копируем наш скрипт автозапуска
COPY /nginx/scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY /nginx/conf/nginx-http.conf /etc/nginx/configs/nginx-http.conf
COPY /nginx/conf/nginx-https.conf /etc/nginx/configs/nginx-https.conf

ENTRYPOINT ["/entrypoint.sh"]