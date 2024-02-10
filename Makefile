# ==================================================================================== #
# VARIABLES
# ==================================================================================== #
# Include variables from the local .env file
include ./app.env

APP_MODULE := github.com/Housiadas/simple-banking-system
DOCKER_COMPOSE_LOCAL := docker-compose -f ./.docker/local/docker-compose.yml
MIGRATE := $(DOCKER_COMPOSE_LOCAL) run --rm utility migrate
SQLC := $(DOCKER_COMPOSE_LOCAL) run --rm utility sqlc
INPUT ?= $(shell bash -c 'read -p "Insert name: " name; echo $$name')

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## docker/build: Build all the containers
.PHONY: docker/build
docker/build:
	$(DOCKER_COMPOSE_LOCAL) build --no-cache --pull

## docker/up: Start all the containers for the app
.PHONY: docker/up
docker/up:
	$(DOCKER_COMPOSE_LOCAL) up -d db redis

## docker/stop: stop all containers
.PHONY: docker/stop
docker/stop:
	$(DOCKER_COMPOSE_LOCAL) stop

## docker/down: stop and remove all containers
.PHONY: docker/down
docker/down:
	$(DOCKER_COMPOSE_LOCAL) down --remove-orphans

## docker/clean: docker clean all
.PHONY: docker/clean
docker/clean:
	docker system prune && \
    docker image prune && \
    docker volume prune

## go/mock/store: Go mock Store interface
.PHONY: go/mock/store
go/mock/store:
	mockgen -package mockdb -destination business/db/mock/store.go $(APP_MODULE)/business/db Store

## go/run: Run main.go locally
.PHONY: go/run
go/run:
	go run app/main.go

# ==================================================================================== #
# DATABASE
# ==================================================================================== #

## db/migrations/create name=$1: create new migration files
.PHONY: db/migrate/create
db/migrate/create:
	$(MIGRATE) create -seq -ext=.sql -dir=./database/migrations $(INPUT)

## db/migrations/up: apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up:
	$(MIGRATE) -path=./database/migrations -database=${MIGRATION_DB_DSN} up

## db/migrations/down: apply all down database migrations (DROP Database)
.PHONY: db/migrate/down
db/migrate/down:
	$(MIGRATE) -path=./database/migrations -database=${MIGRATION_DB_DSN} down

## db/migrations/local/up: apply all up database migrations local command
.PHONY: db/migrate/local/up
db/migrate/local/up:
	go -path=./database/migrations -database=${MIGRATION_DB_DSN} up

## db/sqlc/init: Create an empty sqlc.yaml settings file
.PHONY: db/sqlc/init
db/sqlc/init:
	$(SQLC) init

## db/sqlc/init: Create an empty sqlc.yaml settings file
.PHONY: db/sqlc/generate
db/sqlc/generate:
	$(SQLC) generate

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

# update: update dependeniecs
.PHONY: update
update:
	go get -u ./...
	go mod verify

# vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	go mod verify

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	go mod tidy
	go mod verify
	go fmt ./...
	go vet ./...
	staticcheck ./...
	go test ./... --vet --cover --short --race

# tests: run tests
.PHONY: tests
tests:
	go test ./... -v --cover --short --race

# coverage: Inspect coverage
.PHONY: coverage
coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	open cover.html

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api
