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
run-user-service: build-user-service docker-stop-user-service
	docker compose -f ${DOCKER_COMPOSE_FILE} up user-service-dev -d

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
	docker compose -f ${DOCKER_COMPOSE_FILE} stop user-service-dev
	docker compose -f ${DOCKER_COMPOSE_FILE} ps

clean-user-service-db:
	rm -rf postgres-user-service-data

docker-restart-user-service: docker-stop-user-service docker-start-user-service 

environment-user-service: create-env-file-user-service clean-user-service-db docker-restart-user-service build-user-service

run-migrate-user-service:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t user-service-dev sh -c "./scripts/run_migrate.sh create $(create)"

run-migrate-user-service-up:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t user-service-dev sh -c "./scripts/run_migrate.sh up"

run-migrate-user-service-down:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t user-service-dev sh -c "./scripts/run_migrate.sh down"


run-unit-test-user-service: ## Run unit tests
run-unit-test-user-service: create-env-file-user-service
	@echo "=================="
	@echo "Running unit tests"
	@echo "=================="
	${RUN_IN_DOCKER} user-service-dev sh -c "./scripts/unit_test.sh"

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
run-listing-view-service: build-listing-view-service docker-stop-listing-view-service
	docker compose -f ${DOCKER_COMPOSE_FILE} up listing-view-service-dev -d

run-listing-view-service-consumer: build-listing-view-service docker-stop-listing-view-service-consumer
	docker compose -f ${DOCKER_COMPOSE_FILE} up listing-view-service-consumer-dev -d

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
	docker compose -f ${DOCKER_COMPOSE_FILE} stop listing-view-service-dev
	docker compose -f ${DOCKER_COMPOSE_FILE} ps

docker-stop-listing-view-service-consumer:
	docker compose -f ${DOCKER_COMPOSE_FILE} stop listing-view-service-consumer-dev

clean-listing-view-service-db:
	rm -rf postgres-listing-service-data

docker-restart-listing-view-service: docker-stop-listing-view-service docker-start-listing-view-service

environment-listing-view-service: create-env-file-listing-view-service clean-listing-view-service-db \
	docker-restart-listing-view-service build-listing-view-service

run-migrate-listing-view-service:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t listing-view-service-dev sh -c "./scripts/run_migrate.sh create $(create)"

run-migrate-listing-view-service-up:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t listing-view-service-dev sh -c "./scripts/run_migrate.sh up"


run-migrate-listing-view-service-down:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec -t listing-view-service-dev sh -c "./scripts/run_migrate.sh down"


run-unit-test-listing-view-service: ## Run unit tests
run-unit-test-listing-view-service: create-env-file-listing-view-service
	@echo "=================="
	@echo "Running unit tests"
	@echo "=================="
	${RUN_IN_DOCKER} listing-view-service-dev sh -c "./scripts/unit_test.sh"


#=====================#
#== NATS Server ==#
#=====================#

run-nats-server:
	docker compose -f ${DOCKER_COMPOSE_FILE} up nats-server -d --build --remove-orphans




#=====================#
#== Gateway Service ==#
#=====================#


build-gateway-service: ## Build application
build-gateway-service:
	@echo "==========================="
	@echo "Building binary"
	@echo "==========================="
	${RUN_IN_DOCKER} gateway-service-dev sh -c "go build -mod=vendor -o bin/app cmd/main.go"

run-gateway-service: ## Run HTTP server
run-gateway-service: build-gateway-service docker-stop-gateway-service
	docker compose -f ${DOCKER_COMPOSE_FILE} up gateway-service-dev -d


create-env-file-gateway-service:
	cp gateway-service/.env.sample gateway-service/.env

docker-start-gateway-service:
	@echo "=========================="
	@echo "Starting Docker Containers"
	@echo "=========================="
	docker compose -f ${DOCKER_COMPOSE_FILE} up gateway-service-dev -d --build --remove-orphans
	docker compose -f ${DOCKER_COMPOSE_FILE} ps

docker-stop-gateway-service:
	@echo "=========================="
	@echo "Stopping Docker Containers"
	@echo "=========================="
	docker compose -f ${DOCKER_COMPOSE_FILE} stop gateway-service-dev
	docker compose -f ${DOCKER_COMPOSE_FILE} ps

docker-restart-gateway-service: docker-stop-gateway-service docker-start-gateway-service

environment-gateway-service: create-env-file-gateway-service \
	docker-restart-gateway-service build-gateway-service

run-unit-test-gateway-service: ## Run unit tests
run-unit-test-gateway-service: create-env-file-gateway-service
	@echo "=================="
	@echo "Running unit tests"
	@echo "=================="
	${RUN_IN_DOCKER} gateway-service-dev sh -c "./scripts/unit_test.sh"



environment-all: environment-user-service environment-listing-view-service environment-gateway-service
migrate-all: run-migrate-user-service-up run-migrate-listing-view-service-up
run-all: docker-start-listing-service run-user-service run-listing-view-service \
	run-listing-view-service-consumer run-gateway-service