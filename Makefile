# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

GO_VERSION := 1.22
UID := $(shell id -u)
GID := $(shell id -g)

APP_MODULE := github.com/Housiadas/backend-system
INPUT ?= $(shell bash -c 'read -p "Insert name: " name; echo $$name')
CURRENT_TIME := $(shell date --iso-8601=seconds)
GIT_VERSION := $(shell git describe --always --dirty --tags --long)
LINKER_FLAGS := "-s -X main.buildTime=${CURRENT_TIME} -X main.version=${GIT_VERSION}"

DOCKER_COMPOSE_LOCAL := docker-compose -f ./docker-compose.yml
MIGRATE := $(DOCKER_COMPOSE_LOCAL) run --rm migrate
MIGRATION_DB_DSN := "postgres://housi:secret123@db:5432/housi_db?sslmode=disable"

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## docker/build: Build the application
.PHONY: docker/build
docker/build:
	docker build --target application \
		-t banking-api:local \
		--build-arg GO_VERSION=$(GO_VERSION) \
 		-f .docker/app/Dockerfile .

## docker/up: Start all the containers for the application
.PHONY: docker/up
docker/up:
	$(DOCKER_COMPOSE_LOCAL) up -d

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
	docker system prune -f  && \
    docker image prune -f && \
    docker volume prune -f

## go/mock/store: Go mock Store interface
.PHONY: go/mock/store
go/mock/store:
	mockgen -package mockdb -destination business/db/mock/store.go $(APP_MODULE)/business/db Store

# ==================================================================================== #
# API Application
# ==================================================================================== #

## go/api/run: Run main.go locally
.PHONY: go/api/run
go/api/run:
	go run app/api/main.go

## go/api/build: build the api application
.PHONY: go/api/build
go/api/build:
	cd app & \
	go build -ldflags=${LINKER_FLAGS} -o=./banking-api

# ==================================================================================== #
# CMD Application
# ==================================================================================== #

## go/cmd/build: Build cmd application
.PHONY: go/cmd/build
go/cmd/build:
	go build -o app/cmd/cmd app/cmd/main.go

## go/cmd/seed: Seed db
.PHONY: go/cmd/seed
go/cmd/seed:
	make go/cmd/build
	app/cmd/cmd seed

## go/cmd/useradd: Add user
.PHONY: go/cmd/useradd
go/cmd/useradd:
	make go/cmd/build
	app/cmd/cmd useradd "chris housi" "example@example.com" "1232455477"

## go/cmd/genkey: Generate key
.PHONY: go/cmd/genkey
go/cmd/genkey:
	make go/cmd/build
	app/cmd/cmd genkey

## go/cmd/userevents: User events
.PHONY: go/cmd/userevents
go/cmd/userevents:
	make go/cmd/build
	app/cmd/cmd userevents

# ==================================================================================== #
# DATABASE
# ==================================================================================== #

## db/migrations/create name=$1: create new migration files
.PHONY: db/migrate/create
db/migrate/create:
	$(MIGRATE) create -seq -ext=.sql -dir=./business/data/migrations $(INPUT)

## db/migrations/up: apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up:
	$(MIGRATE) -path=./business/data/migrations -database=${MIGRATION_DB_DSN} up

## db/migrations/down: apply all down database migrations (DROP Database)
.PHONY: db/migrate/down
db/migrate/down:
	$(MIGRATE) -path=./business/data/migrations -database=${MIGRATION_DB_DSN} down

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

# vendor: vendor dependencies
.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	go mod verify

# update: update dependencies
.PHONY: update
update:
	go get -u ./...
	go mod verify

## lint: Go linter
.PHONY: lint
lint:
	go mod tidy
	go mod verify
	go fmt ./...
	go vet ./...
	staticcheck ./...

# tests: run tests
.PHONY: tests
tests:
	go install github.com/mfridman/tparse@latest
	CGO_ENABLED=1 go test -v --cover --short --race -json ./... | tparse --all

# coverage: Inspect coverage
.PHONY: coverage
coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	open cover.html

# ==================================================================================== #
# SWAGGER
# ==================================================================================== #

# swagger: Generate swagger docs
.PHONY: swagger
swagger:
	docker run --rm -v $(PWD):/code --user $(UID) ghcr.io/swaggo/swag:v1.16.3 init --g app/api/main.go
