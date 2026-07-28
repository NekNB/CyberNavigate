function Show-Status {
    Clear-Host
    docker stack ps -f desired-state=running -f desired-state=ready cyber-navigate
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