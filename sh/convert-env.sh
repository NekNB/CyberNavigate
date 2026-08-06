#!/bin/bash

# Получаем директорию проекта (на два уровня выше текущего скрипта)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
workspaceFolder="$(dirname "$SCRIPT_DIR")"

# Ищем все файлы *.tmp.env рекурсивно в папке secrets/templates
# -print0 и read -d '' используются для безопасной обработки путей с пробелами
find "$workspaceFolder/secrets/templates" -type f -name "*.tmp.env" -print0 | while IFS= read -r -d '' file; do

    # Заменяем .tmp.env -> .env (отрезаем суффикс и добавляем .env)
    newFilePath="${file%.tmp.env}.env"

    # Заменяем /templates/ -> /env/ (глобальная замена в строке)
    # Экранируем слэши, чтобы не конфликтовать с разделителем replace
    newFilePath="${newFilePath//\/templates\//\/env\/}"

    # Создаем директорию для файла, если её нет
    newDir="$(dirname "$newFilePath")"
    mkdir -p "$newDir"

    # Копируем файл
    cp "$file" "$newFilePath"

    echo "Создан: $newFilePath"
done

echo "Готово"