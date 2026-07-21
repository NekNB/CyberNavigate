# Получаем директорию проекта
$workspaceFolder = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

# Ищем все файлы *.tmp.env рекурсивно в папке secrets/templates
Get-ChildItem -Path $workspaceFolder/secrets/templates -Recurse -Filter "*.tmp.env" | ForEach-Object {

    # Заменяем .tmp.env -> .env и templates -> env
    $newFilePath = $_.FullName -replace '\.tmp\.env$', '.env'
    $newFilePath = $newFilePath -replace '\\templates\\', '\\env\\'

    # Создаем директорию для файла, если её нет
    $newDir = Split-Path $newFilePath -Parent
    if (!(Test-Path $newDir)) {
        New-Item -ItemType Directory -Path $newDir -Force | Out-Null
    }

    # Копируем содержимое в новый файл
    Copy-Item -Path $_.FullName -Destination $newFilePath -Force

    Write-Host "Создан: $newFilePath"
}

Write-Host "Готово"