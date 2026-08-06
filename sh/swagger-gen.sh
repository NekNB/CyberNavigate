#!/bin/bash

# Проверка наличия обязательного аргумента
if [[ -z "$1" ]]; then
    echo -e "\033[31mОшибка: Не указан параметр FileName\033[0m"
    echo "Использование: $0 <FileName>"
    exit 1
fi

FileName="$1"

# Функция замены строки в файле
set_file_content_replacement() {
    local FilePath="$1"
    local OldString="$2"
    local NewString="$3"
    
    # Проверяем существование файла
    if [[ ! -f "$FilePath" ]]; then
        echo -e "\033[31mОшибка: Файл '$FilePath' не найден\033[0m" >&2
        return 1
    fi
    
    # Читаем содержимое файла в переменную
    local content
    content=$(cat "$FilePath")
    
    # Проверяем наличие старой строки (grep -qF ищет точное совпадение литералов)
    if ! grep -qF -- "$OldString" "$FilePath"; then
        echo -e "\033[33mПредупреждение: Строка не найдена в файле: $FilePath\033[0m" >&2
        return 1
    fi
    
    # Заменяем строку (нативный метод Bash работает с литералами, экранируя спецсимволы)
    local newContent="${content//"$OldString"/$NewString}"
    
    # Сохраняем изменения без добавления переноса строки в конце (-NoNewline)
    printf '%s' "$newContent" > "$FilePath"
    
    if [[ $? -eq 0 ]]; then
        echo -e "\033[32m✓ Замена выполнена в файле: $FilePath\033[0m"
        return 0
    else
        echo -e "\033[31mОшибка при записи в файл '$FilePath'\033[0m" >&2
        return 1
    fi
}

# Switch-case по имени файла
case "$FileName" in
    "user.swagger.yaml")
        echo -e "\033[36mОбработка файла: user.swagger.yaml\033[0m"
        pushd ../swagger/ > /dev/null || exit 1
        
        oapi-codegen -config ./configs/user/models.yaml ./docs/user-service/user.swagger.yaml
        oapi-codegen -config ./configs/user/client.yaml ./docs/user-service/user.swagger.yaml
        oapi-codegen -config ./configs/user/server.yaml ./docs/user-service/user.swagger.yaml
        
        set_file_content_replacement "./gen/user/server.go" \
            '*fiber.Ctx' \
            'fiber.Ctx'

        set_file_content_replacement "./gen/user/server.go" \
            'github.com/gofiber/fiber/v2' \
            'github.com/gofiber/fiber/v3'

        set_file_content_replacement "./gen/user/server.go" \
            'c.Context().SetUserValue((BearerAuthScopes), []string{})' \
            'c.Locals(BearerAuthScopes, []string{})'
            
        popd > /dev/null
        ;;
        
    "article.swagger.yaml")
        echo -e "\033[36mОбработка файла: article.swagger.yaml\033[0m"
        pushd ../swagger/ > /dev/null || exit 1
        
        oapi-codegen -config ./configs/article/models.yaml ./docs/article-service/article.swagger.yaml
        oapi-codegen -config ./configs/article/client.yaml ./docs/article-service/article.swagger.yaml
        oapi-codegen -config ./configs/article/server.yaml ./docs/article-service/article.swagger.yaml
        
        set_file_content_replacement "./gen/article/server.go" \
            '*fiber.Ctx' \
            'fiber.Ctx'

        set_file_content_replacement "./gen/article/server.go" \
            'github.com/gofiber/fiber/v2' \
            'github.com/gofiber/fiber/v3'

        popd > /dev/null
        ;;
        
    "simulator.swagger.yaml")
        echo -e "\033[36mОбработка файла: simulator.swagger.yaml\033[0m"
        pushd ../swagger/ > /dev/null || exit 1
        
        oapi-codegen -config ./configs/simulator/models.yaml ./docs/simulator-service/simulator.swagger.yaml
        oapi-codegen -config ./configs/simulator/client.yaml ./docs/simulator-service/simulator.swagger.yaml
        oapi-codegen -config ./configs/simulator/server.yaml ./docs/simulator-service/simulator.swagger.yaml
        
        set_file_content_replacement "./gen/simulator/server.go" \
            '*fiber.Ctx' \
            'fiber.Ctx'
            
        set_file_content_replacement "./gen/simulator/server.go" \
            'github.com/gofiber/fiber/v2' \
            'github.com/gofiber/fiber/v3'
            
        set_file_content_replacement "./gen/simulator/server.go" \
            'c.Context().SetUserValue((BearerAuthScopes), []string{})' \
            'c.Locals(BearerAuthScopes, []string{})'
            
        popd > /dev/null
        ;;
        
    *)
        echo -e "\033[31mНеизвестный файл: $FileName\033[0m"
        exit 1
        ;;
esac