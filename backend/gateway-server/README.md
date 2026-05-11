env - не требуется

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
