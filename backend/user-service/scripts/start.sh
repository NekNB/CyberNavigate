# !/bin/sh

SECRET_PATH="/run/secrets/user_service_secret"

load_secrets() {
	if [ -f "$SECRET_PATH" ]; then
			# Используем set -a чтобы экспортировать все переменные из файла автоматически
			# Это работает, если файл содержит VAR=Value строки
			set -a
			source <(tr -d '\r' < $SECRET_PATH)
			set +a
			
			# Проверка, загрузились ли переменные	
			if [ -z "$POSTGRES_PASSWORD" ] || [ -z "$MONGO_PASSWORD" ]; then
					echo "Error: Failed to load secrets from $SECRET_PATH or variables are empty."
					exit 1
			fi
	else
			echo "Warning: Secret file not found at $SECRET_PATH"
	fi
}


load_secrets
./user-service