.PHONY: echorun build start stop secrets_remove srm secrets_create scr secrets_update


.SILENT:

ARGS = $(filter-out $@,$(MAKECMDGOALS))

run:
	$(MAKE) build ARGS=postgres
	$(MAKE) build ARGS=mongo
	$(MAKE) build ARGS=article-service
	$(MAKE) build ARGS=gateway-server
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

secrets_create scr:
	docker secret create postgres_secret ./secrets/postgres.secret.env || echo ">> Secret already exists"
	docker secret create mongo_secret ./secrets/mongo.secret.env || echo ">> Secret already exists"
	docker secret create article_service_secret ./secrets/article.secret.env || echo ">> Secret already exists"
	

start:
	docker stack deploy -c docker-compose.yaml cyber-navigate

scale:
	$(eval SERVICE=$(word 1,$(ARGS)))
	$(eval COUNT=$(word 2,$(ARGS)))
	docker service scale cyber-navigate_$(SERVICE)=$(COUNT)

stop:
	docker service scale cyber-navigate_${ARGS}=0

down:
	docker stack rm cyber-navigate


update:
	docker service update --force cyber-navigate_${ARGS}


%:
	@: