# Проверка наличия обязательного аргумента
def main [FileName: string] {
    match $FileName {
        "user.swagger.yaml" => {
            print "\u{1b}[36mОбработка файла: user.swagger.yaml\u{1b}[0m"
            process-swagger "user" [
                {
                    old: '*fiber.Ctx'
                    new: 'fiber.Ctx'
                }
                {
                    old: 'github.com/gofiber/fiber/v2'
                    new: 'github.com/gofiber/fiber/v3'
                }
                {
                    old: 'c.Context().SetUserValue((BearerAuthScopes), []string{})'
                    new: 'c.Locals(BearerAuthScopes, []string{})'
                }
            ]
        }

        "article.swagger.yaml" => {
            print "\u{1b}[36mОбработка файла: article.swagger.yaml\u{1b}[0m"

            process-swagger "article" [
                {
                    old: '*fiber.Ctx'
                    new: 'fiber.Ctx'
                }
                {
                    old: 'github.com/gofiber/fiber/v2'
                    new: 'github.com/gofiber/fiber/v3'
                }
            ]
        }

        "simulator.swagger.yaml" => {
            print "\u{1b}[36mОбработка файла: simulator.swagger.yaml\u{1b}[0m"

            process-swagger "simulator" [
                {
                    old: '*fiber.Ctx'
                    new: 'fiber.Ctx'
                }
                {
                    old: 'github.com/gofiber/fiber/v2'
                    new: 'github.com/gofiber/fiber/v3'
                }
                {
                    old: 'c.Context().SetUserValue((BearerAuthScopes), []string{})'
                    new: 'c.Locals(BearerAuthScopes, []string{})'
                }
            ]
        }

        _ => {
            print --stderr $"\u{1b}[31mНеизвестный файл: ($FileName)\u{1b}[0m"
            exit 1
        }
    }
}

# Функция замены строки в файле
def set-file-content-replacement [
    FilePath: string
    OldString: string
    NewString: string
] {
    if not ($FilePath | path exists) {
        print --stderr $"\u{1b}[31mОшибка: Файл '($FilePath)' не найден\u{1b}[0m"
        return false
    }

    let content = open --raw $FilePath

    # Проверяем наличие старой строки
    if not ($content | str contains -- $OldString) {
        print --stderr $"\u{1b}[33mПредупреждение: Строка не найдена в файле: ($FilePath)\u{1b}[0m"
        return false
    }

    # Литеральная замена всех вхождений
    let new_content = $content | str replace --all -- $OldString $NewString

    # Запись без добавления переноса строки
    try {
        $new_content | save --force --raw $FilePath

        print $"\u{1b}[32m✓ Замена выполнена в файле: ($FilePath)\u{1b}[0m"
        return true
    } catch {
        print --stderr $"\u{1b}[31mОшибка при записи в файл '($FilePath)'\u{1b}[0m"
        return false
    }
}

# Генерация OpenAPI-кода
def process-swagger [
    service: string
    replacements?: list<record<old: string, new: string>>
] {
    let original_dir = $env.PWD
    let swagger_dir = "../swagger"

    cd $swagger_dir

    try {
        # Генерация моделей, клиента и сервера
        for config in ["models" "client" "server"] {
            let config_path = $"./configs/($service)/($config).yaml"
            let swagger_path = $"./docs/($service)-service/($service).swagger.yaml"

            ^oapi-codegen -config $config_path $swagger_path

            if $env.LAST_EXIT_CODE != 0 {
                error make { msg: $"Ошибка генерации: ($config_path)" }
            }
        }

        # Замена строк в сгенерированном сервере
        let server_file = $"./gen/($service)/server.go"
        if (not ($replacements | is-empty)) {
            for replacement in $replacements {
                let result = set-file-content-replacement $server_file $replacement.old $replacement.new
                
                if not $result {
                    print --stderr $"Предупреждение: Не удалось выполнить замену ($replacement.old) на ($replacement.new)"
                }
            }
        }
    } catch {|err|
        print --stderr $"Ошибка обработки сервиса ($service): ($err.msg)"
        cd $original_dir
        exit 1
    }

    cd $original_dir
}