DOCKER_COMPOSE_FILE=docker-compose.yaml

.PHONY: build up start stop down logs ps restart prune


build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

start:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d --build

stop:
	docker-compose -f $(DOCKER_COMPOSE_FILE) stop

down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

ps:
	docker-compose -f $(DOCKER_COMPOSE_FILE) ps

restart:
	docker-compose -f $(DOCKER_COMPOSE_FILE) restart

prune:
	docker system prune -af --volumes
