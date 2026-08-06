#!/bin/bash

# Получаем директорию проекта (на два уровня выше текущего скрипта)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
workspaceFolder="$(dirname "$SCRIPT_DIR")"

# Запоминаем текущую директорию и переходим в secrets
pushd "$workspaceFolder/secrets" > /dev/null || exit 1

output="Ключи уже существуют"

# Создаем папку keys, если её нет
if [[ ! -d ./keys ]]; then
    mkdir -p ./keys
fi

# Если нет приватного ключа, генерируем новые
if [[ ! -f ./keys/private.pem ]]; then
    cd ./keys || exit 1
    openssl genrsa -out private.pem 2048
    openssl rsa -in private.pem -pubout -out public.pem
    
    # Присваиваем новое значение переменной (в PS тут была ошибка с записью в файл)
    output="Созданы новые ключи"
fi

echo "$output"

# Возвращаемся в исходную директорию
popd > /dev/null || exit 1