set shell := ["nu", "-c"]

build env service:
	docker build -f ./docker/{{ env }}/{{ service }}.dockerfile -t cyber-navigate/{{ service }} .

run env: 
	just build {{ env }} nginx
	just build {{ env }} postgres
	just build {{ env }} mongo
	just build {{ env }} article-service
	just build {{ env }} user-service
	just build {{ env }} simulator-service
	just build {{ env }} gateway-server
	just build {{ env }} frontend

	just secrets_update
	just deploy

deploy:
	docker stack deploy -c docker-compose.yaml cyber-navigate

[no-exit-message]
secrets_remove: 
	-docker secret remove postgres_secret 
	-docker secret remove mongo_secret 
	-docker secret remove article_service_secret 
	-docker secret remove user_service_secret 
	-docker secret remove simulator_service_secret 
[no-exit-message]
secrets_create:
	-docker secret create postgres_secret ./secrets/env/postgres.secret.env 
	-docker secret create mongo_secret ./secrets/env/mongo.secret.env 
	-docker secret create article_service_secret ./secrets/env/article.secret.env 
	-docker secret create user_service_secret ./secrets/env/user.secret.env 
	-docker secret create simulator_service_secret ./secrets/env/simulator.secret.env 


secrets_update: secrets_remove secrets_create

start service:
	docker service scale cyber-navigate_{{ service }}=1

stop service:
	docker service scale cyber-navigate_{{ service }}=0

update service:
	docker service update --force cyber-navigate_{{ service }}

down:
	docker stack rm cyber-navigate
