# Article Service

Работает с Mongo и Postgres
В Mongo - хранит текст статей
В Postgres - хранит сущности

Swagger доступен на /swagger

secrets/article.secret.env

```env
POSTGRES_PASSWORD=article_service_password
MONGO_PASSWORD=article_password
```

configs/gateway-server/dev.yaml

```yaml
env: "dev" # Тип переменного окружения local/dev/prod

server:
  port: 9000 # Port запуска. Прописан в Docker Compose
  services: # Перечисление сервисов
    - article:
        path: /api/v1/articles
        protocol: http
        host: cyber-navigate_article-service
        port: 8000
```

Запуск: `go run ./cmd --config=config_path`
