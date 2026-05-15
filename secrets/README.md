# Ключи

Приватный: `openssl genrsa -out private.pem 2048`

Публичный: `openssl rsa -in private.pem -pubout -out public.pem`
