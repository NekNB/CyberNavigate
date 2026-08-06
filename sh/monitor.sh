#!/bin/bash

show_status() {
    clear
    docker stack ps -f desired-state=running -f desired-state=ready cyber-navigate
    echo ""
    echo "R — обновить"
    echo "Q — выйти"
}

# Первичный вызов функции
show_status

while true; do
    # Читаем одну нажатую клавишу без ожидания Enter
    read -rsn1 key

    # Проверяем нажатую клавишу (с учетом регистра)
    if [[ "$key" == "r" || "$key" == "R" ]]; then
        show_status
    elif [[ "$key" == "q" || "$key" == "Q" ]]; then
        break
    fi
done