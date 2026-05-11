Генерация Go Файлов:

Конфигурация OApi-Codegen: ./configs/{serviceName}/...

Для Article:
oapi-codegen -config ./configs/article/models.yaml ./docs/article-service/article.swagger.yaml
oapi-codegen -config ./configs/article/client.yaml ./docs/article-service/article.swagger.yaml
oapi-codegen -config ./configs/article/server.yaml ./docs/article-service/article.swagger.yaml