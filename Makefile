.PHONY: run build start stop secrets_remove srm secrets_create scr secrets_update

run: stop build secrets_update start

build:
	docker build -f ./postgres/Dockerfile -t cyber-navigate/postgres .

secrets_update: secrets_remove secrets_create

secrets_remove srm:
	docker secret remove postgres_secret || echo ">> Secret not found"

secrets_create scr:
	docker secret create postgres_secret ./secrets/postgres.secret




start:
	docker stack deploy -c docker-compose.yaml cyber-navigate


stop:
	docker stack rm cyber-navigate

