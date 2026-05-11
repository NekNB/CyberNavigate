# Получаем директорию, где лежит сам скрипт
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Ищем все файлы *.tmp.env рекурсивно
Get-ChildItem -Path $scriptDir -Recurse -Filter "*.tmp.env" | ForEach-Object {

    # Новый путь: заменяем .tmp.env -> .env
    $newFilePath = $_.FullName -replace '\.tmp\.env$', '.env'

    # Копируем содержимое в новый файл
    Copy-Item -Path $_.FullName -Destination $newFilePath -Force

    Write-Host "Создан: $newFilePath"
}

Write-Host "Готово"