.PHONY: build start stop secrets_remove srm secrets_create scr secrets_update ps1 sh

.SILENT:

ARGS = $(filter-out $@,$(MAKECMDGOALS))

ps1:
	powershell -Command "Push-Location ./ps1; ./convert-env.ps1; ./create-keys.ps1; ./generate-mtls-certs.ps1; Pop-Location"


sh:
	chmod +x ./sh/convert-env.sh ./sh/create-keys.sh ./sh/generate-mtls-certs.sh
	cd ./sh && ./convert-env.sh && ./create-keys.sh && ./generate-mtls-certs.sh 

run:
	$(MAKE) build ARGS=nginx
	$(MAKE) build ARGS=postgres
	$(MAKE) build ARGS=mongo
	$(MAKE) build ARGS=article-service
	$(MAKE) build ARGS=user-service
	$(MAKE) build ARGS=simulator-service
	$(MAKE) build ARGS=gateway-server
	$(MAKE) build ARGS=frontend
	$(MAKE) secrets_update
	$(MAKE) start


# Требует ввести name
build:
	docker build -f ./docker/${ARGS}.dockerfile -t cyber-navigate/${ARGS} .




secrets_update: secrets_remove secrets_create

secrets_remove srm:
	docker secret remove postgres_secret || echo ">> Secret not found"
	docker secret remove mongo_secret || echo ">> Secret not found"
	docker secret remove article_service_secret || echo ">> Secret not found"
	docker secret remove user_service_secret || echo ">> Secret not found"
	docker secret remove simulator_service_secret || echo ">> Secret not found"

secrets_create scr:
	docker secret create postgres_secret ./secrets/env/postgres.secret.env || echo ">> Secret already exists"
	docker secret create mongo_secret ./secrets/env/mongo.secret.env || echo ">> Secret already exists"
	docker secret create article_service_secret ./secrets/env/article.secret.env || echo ">> Secret already exists"
	docker secret create user_service_secret ./secrets/env/user.secret.env || echo ">> Secret already exists"
	docker secret create simulator_service_secret ./secrets/env/simulator.secret.env || echo ">> Secret already exists"

start:
	docker stack deploy -c docker-compose.yaml cyber-navigate

scale:
	$(eval SERVICE=$(word 1,$(ARGS)))
	$(eval COUNT=$(word 2,$(ARGS)))
	docker service scale cyber-navigate_$(SERVICE)=$(COUNT)
start_service:
	docker service scale cyber-navigate_$(ARGS)=1

stop:
	docker service scale cyber-navigate_${ARGS}=0

down:
	docker stack rm cyber-navigate


update:
	docker service update --force cyber-navigate_${ARGS}


%:
	@: