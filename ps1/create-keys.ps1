# Получаем директорию проекта
$workspaceFolder =  Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)

Push-Location $workspaceFolder/secrets 

$output = "Ключи уже существуют"

$null = & {
	if (!(Test-Path ./keys)) {mkdir keys}

	if (!(Test-Path ./keys/private.pem)) {
		Set-Location ./keys
		openssl genrsa -out private.pem 2048 
		openssl rsa -in private.pem -pubout -out public.pem
		"Созданы новые ключи" > $output
	}
}

Write-Output $output

Pop-Location 