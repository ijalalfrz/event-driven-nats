DOCKER_COMPOSE_FILE=docker-compose.yml
RUN_IN_DOCKER ?= docker compose -f ${DOCKER_COMPOSE_FILE} exec -T
WITH_TTY ?= -t

#=====================#
#== Listing Service ==#
#=====================#

docker-start-listing-service:
	docker compose -f $(DOCKER_COMPOSE_FILE) up listing-service -d --build --remove-orphans

docker-stop-listing-service:
	docker compose -f $(DOCKER_COMPOSE_FILE) stop listing-service

docker-restart-listing-service: docker-stop-listing-service docker-start-listing-service


#===================#
#== User Service ==#
#===================#

build-user-service: ## Build application
build-user-service:
	@echo "==========================="
	@echo "Building binary"
	@echo "==========================="
	${RUN_IN_DOCKER} user-service-dev sh -c "go build -mod=vendor -o bin/app cmd/main.go"

run-user-service: ## Run HTTP server
run-user-service: build-user-service
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t user-service-dev sh -c ./scripts/run_server.sh


create-env-file-user-service:
	cp user-service/.env.sample user-service/.env

docker-start-user-service:
	@echo "=========================="
	@echo "Starting Docker Containers"
	@echo "=========================="
	docker compose -f ${DOCKER_COMPOSE_FILE} up user-service-dev postgres-user-service -d --build --remove-orphans
	docker compose -f ${DOCKER_COMPOSE_FILE} ps

docker-stop-user-service:
	@echo "=========================="
	@echo "Stopping Docker Containers"
	@echo "=========================="
	docker compose -f ${DOCKER_COMPOSE_FILE} stop user-service-dev postgres-user-service
	docker compose -f ${DOCKER_COMPOSE_FILE} ps

docker-restart-user-service: docker-stop-user-service docker-start-user-service

environment-user-service: create-env-file-user-service docker-restart-user-service build-user-service


#=====================#
#== Listing View Service ==#
#=====================#


build-listing-view-service: ## Build application
build-listing-view-service:
	@echo "==========================="
	@echo "Building binary"
	@echo "==========================="
	${RUN_IN_DOCKER} listing-view-service-dev sh -c "go build -mod=vendor -o bin/app cmd/main.go"

run-listing-view-service: ## Run HTTP server
run-listing-view-service: build-listing-view-service
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t listing-view-service-dev sh -c ./scripts/run_server.sh


create-env-file-listing-view-service:
	cp listing-view-service/.env.sample listing-view-service/.env

docker-start-listing-view-service:
	@echo "=========================="
	@echo "Starting Docker Containers"
	@echo "=========================="
	docker compose -f ${DOCKER_COMPOSE_FILE} up listing-view-service-dev postgres-listing-service -d --build --remove-orphans
	docker compose -f ${DOCKER_COMPOSE_FILE} ps

docker-stop-listing-view-service:
	@echo "=========================="
	@echo "Stopping Docker Containers"
	@echo "=========================="
	docker compose -f ${DOCKER_COMPOSE_FILE} stop listing-view-service-dev postgres-listing-service
	docker compose -f ${DOCKER_COMPOSE_FILE} ps

docker-restart-listing-view-service: docker-stop-listing-view-service docker-start-listing-view-service

environment-listing-view-service: create-env-file-listing-view-service docker-restart-listing-view-service build-listing-view-service

run-listing-view-service-consumer: build-listing-view-service
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t listing-view-service-dev sh -c ./scripts/run_consumer.sh


run-nats-server:
	docker compose -f ${DOCKER_COMPOSE_FILE} up nats-server -d --build --remove-orphans

