set shell := ["nu", "-c"]

build env service:
	docker build -f ./docker/{{ env }}/{{ service }}.dockerfile -t cyber-navigator/{{ service }} .

tag service:
	docker tag cyber-navigator/{{ service }}:latest neknarej/cyber-navigator_{{ service }}:latest	

push service:
	docker push neknarej/cyber-navigator_{{ service }}:latest
send:
	just build prod nginx
	just tag nginx
	just push nginx
	 

	just build prod postgres
	just tag postgres
	just push postgres
	 

	just build prod mongo
	just tag mongo 
	just push mongo
	 

	just build prod article-service
	just tag article-service
	just push article-service
	 

	just build prod user-service
	just tag user-service
	just push user-service
	 

	just build prod simulator-service
	just tag simulator-service
	just push simulator-service
	 

	just build prod gateway-server
	just tag gateway-server
	just push gateway-server
	 

	just build prod frontend
	just tag frontend
	just push frontend
	 



run env: 
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
