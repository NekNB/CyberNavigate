# simulator Service

Работает с Postgres


Swagger доступен на /swagger

secrets/simulator.secret.env

```env
POSTGRES_PASSWORD=simulator_service_password
```

configs/gateway-server/dev.yaml

```yaml
env: "dev" # Тип переменного окружения local/dev/prod

server:
  port: 9000 # Port запуска. Прописан в Docker Compose
  services: # Перечисление сервисов
    - simulator:
        path: /api/v1/simulators
        protocol: http
        host: cyber-navigate_simulator-service
        port: 8000
```

Запуск: `go run ./cmd --config=config_path`

Некоторые ручки требуют установку X-User-Id в Header
