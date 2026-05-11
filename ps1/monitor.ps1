function Show-Status {
    Clear-Host
    docker stack ps cyber-navigate
    Write-Host ""
    Write-Host "R — обновить"
		 Write-Host "Q — выйти"
}

Show-Status

while ($true) {
    $key = [System.Console]::ReadKey($true)

    if ($key.Key -eq [ConsoleKey]::R) {
        Show-Status
    }
		elseif ($key.Key -eq [ConsoleKey]::Q) {
        break
    }
}