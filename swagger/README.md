Генерация Go Файлов:

Конфигурация OApi-Codegen: ./configs/{serviceName}/...

Для Article:
oapi-codegen -config ./configs/article/models.yaml ./docs/article-service/article.swagger.yaml
oapi-codegen -config ./configs/article/client.yaml ./docs/article-service/article.swagger.yaml
oapi-codegen -config ./configs/article/server.yaml ./docs/article-service/article.swagger.yaml

Для User:
oapi-codegen -config ./configs/user/models.yaml ./docs/user-service/user.swagger.yaml
oapi-codegen -config ./configs/user/client.yaml ./docs/user-service/user.swagger.yaml
oapi-codegen -config ./configs/user/server.yaml ./docs/user-service/user.swagger.yaml

    c.Context().SetUserValue((BearerAuthScopes), []string{}) >> c.Locals(BearerAuthScopes, []string{})
