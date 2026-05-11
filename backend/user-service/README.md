# User Service

Работает c Postgres
В Postgres - хранит сущности

Swagger доступен на /swagger

secrets/user.secret.env

```env
POSTGRES_PASSWORD=user_service_password
MONGO_PASSWORD=user_password
```

configs/gateway-server/dev.yaml

```yaml
env: "dev" # Тип переменного окружения local/dev/prod

server:
  port: 9000 # Port запуска. Прописан в Docker Compose
  services: # Перечисление сервисов
    - user:
        path: /api/v1/users
        protocol: http
        host: cyber-navigate_user-service
        port: 8000
```

Запуск: `go run ./cmd --config=config_path`
