param(
    [Parameter(Mandatory=$true)]
    [string]$FileName
)

# Функция замены строки в файле
function  Set-FileContentReplacement {
    param(
        [Parameter(Mandatory=$true)]
        [string]$FilePath,
        
        [Parameter(Mandatory=$true)]
        [string]$OldString,
        
        [Parameter(Mandatory=$true)]
        [string]$NewString
    )
    
    # Проверяем существование файла
    if (-not (Test-Path $FilePath)) {
        Write-Error "Файл '$FilePath' не найден"
        return $false
    }
    
    try {
        # Читаем содержимое файла
        $content = Get-Content -Path $FilePath -Raw -ErrorAction Stop
        
        # Проверяем наличие старой строки
        if ($content -notmatch [regex]::Escape($OldString)) {
            Write-Warning "Строка не найдена в файле: $FilePath"
            return $false
        }
        
        # Заменяем строку
        $newContent = $content -replace [regex]::Escape($OldString), $NewString
        
     
        
        # Сохраняем изменения
        Set-Content -Path $FilePath -Value $newContent -NoNewline -ErrorAction Stop
        
        Write-Host "✓ Замена выполнена в файле: $FilePath" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Error "Ошибка при обработке файла '$FilePath': $_"
        return $false
    }
}

# Switch-case по имени файла
switch ($FileName) {
    "user.swagger.yaml" {
        Write-Host "Обработка файла: user.swagger.yaml" -ForegroundColor Cyan
        Push-Location ../swagger/
   
        oapi-codegen -config ./configs/user/models.yaml ./docs/user-service/user.swagger.yaml
        oapi-codegen -config ./configs/user/client.yaml ./docs/user-service/user.swagger.yaml
        oapi-codegen -config ./configs/user/server.yaml ./docs/user-service/user.swagger.yaml
        
        Set-FileContentReplacement -FilePath ./gen/user/server.go `
            -OldString '*fiber.Ctx' `
            -NewString 'fiber.Ctx'

        Set-FileContentReplacement -FilePath ./gen/user/server.go `
            -OldString 'github.com/gofiber/fiber/v2' `
            -NewString 'github.com/gofiber/fiber/v3'

        Set-FileContentReplacement -FilePath ./gen/user/server.go `
            -OldString 'c.Context().SetUserValue((BearerAuthScopes), []string{})' `
            -NewString 'c.Locals(BearerAuthScopes, []string{})'
        Pop-Location
        break
    }
    
    "article.swagger.yaml" {
        Write-Host "Обработка файла: article.swagger.yaml" -ForegroundColor Cyan
        Push-Location ../swagger/
        oapi-codegen -config ./configs/article/models.yaml ./docs/article-service/article.swagger.yaml
        oapi-codegen -config ./configs/article/client.yaml ./docs/article-service/article.swagger.yaml
        oapi-codegen -config ./configs/article/server.yaml ./docs/article-service/article.swagger.yaml
        
        Set-FileContentReplacement -FilePath ./gen/article/server.go `
            -OldString '*fiber.Ctx' `
            -NewString 'fiber.Ctx'

        Set-FileContentReplacement -FilePath ./gen/article/server.go `
            -OldString 'github.com/gofiber/fiber/v2' `
            -NewString 'github.com/gofiber/fiber/v3'

        Pop-Location
        break
    }
    
   
    
    default {
        Write-Host "Неизвестный файл: $FileName" -ForegroundColor Red
        exit 1
    }
}