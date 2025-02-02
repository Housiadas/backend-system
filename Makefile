# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

GO_VERSION := 1.23
UID := $(shell id -u)
GID := $(shell id -g)

APP_MODULE := github.com/Housiadas/backend-system
INPUT ?= $(shell bash -c 'read -p "Insert name: " name; echo $$name')
CURRENT_TIME := $(shell date --iso-8601=seconds)
GIT_VERSION := $(shell git describe --always --dirty --tags --long)
LINKER_FLAGS := "-s -X main.buildTime=${CURRENT_TIME} -X main.version=${GIT_VERSION}"

DOCKER_COMPOSE_LOCAL := docker compose -f ./docker-compose.yml
MIGRATE := $(DOCKER_COMPOSE_LOCAL) run --rm migrate
MIGRATION_DB_DSN := "postgres://housi:secret123@db:5432/housi_db?sslmode=disable"

## ========
## Docker
## ========

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

## ==================
## HTTP Application
## ==================

## go/http/run: Run main.go locally
.PHONY: go/api/run
go/http/run:
	go run app/api/main.go

## go/http/build: build the http application
.PHONY: go/api/build
go/http/build:
	cd app & \
	go build -ldflags=${LINKER_FLAGS} -o=./banking-api

## ==================
## CMD Application
## ==================

## go/cmd/build: Build cmd application
.PHONY: go/cmd/build
go/cmd/build:
	go build -o app/cmd/cmd app/cmd/main.go

## go/cmd/seed: Seed db
.PHONY: go/cmd/seed
go/cmd/seed:
	make go/cmd/build
	app/cmd/cmd seed

## go/cmd/genkey: Generate key
.PHONY: go/cmd/genkey
go/cmd/genkey:
	make go/cmd/build
	app/cmd/cmd genkey

## go/cmd/useradd: Add user
.PHONY: go/cmd/useradd
go/cmd/user/add:
	make go/cmd/build
	app/cmd/cmd useradd "chris housi" "example@example.com" "1232455477"

## go/cmd/user/events: User events
.PHONY: go/cmd/userevents
go/cmd/user/events:
	make go/cmd/build
	app/cmd/cmd userevents

## ==================
## Database
## ==================

## db/migrations/create name=$1: Create new migration files
.PHONY: db/migrate/create
db/migrate/create:
	$(MIGRATE) create -seq -ext=.sql -dir=./business/data/migrations $(INPUT)

## db/migrations/up: Apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up:
	$(MIGRATE) -path=./business/data/migrations -database=${MIGRATION_DB_DSN} up

## db/migrations/down: Apply all down database migrations (DROP Database)
.PHONY: db/migrate/down
db/migrate/down:
	$(MIGRATE) -path=./business/data/migrations -database=${MIGRATION_DB_DSN} down

## ==================
## Quality Control
## ==================

## lint: Run linter
.PHONY: lint
lint:
	go mod tidy
	go mod verify
	go fmt ./...
	go vet ./...
	staticcheck ./...

## tests: Run tests
.PHONY: tests
tests:
	go install github.com/mfridman/tparse@latest
	CGO_ENABLED=1 go test -v --cover --short --race -json ./... | tparse --all

## coverage: Inspect coverage
.PHONY: coverage
coverage:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	open cover.html

## ==================
## Modules support
## ==================

## deps/vendor: Vendor dependencies
.PHONY: vendor
deps/vendor:
	go mod tidy
	go mod vendor
	go mod verify

## deps/update: Update dependencies
.PHONY: deps/update
deps/update:
	go get -u -v ./...
	go mod tidy
	go mod vendor

## deps/list: List dependencies
.PHONY: deps/list
deps/list:
	go list -m -u -mod=readonly all

## deps/cache/clean: Clean cache dependencies
.PHONY: deps/cache/clean
deps/cache/clean:
	go clean -modcache

## deps/reset: Reset dependencies
.PHONY: deps/reset
deps/reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

## list: List modules
.PHONY: list
list:
	go list -mod=mod all

## ==================
## Utils
## ==================

## go/mock/store: Go mock Store interface
.PHONY: go/mock/store
go/mock/store:
	mockgen -package mockdb -destination business/db/mock/store.go $(APP_MODULE)/business/db Store

# swagger: Generate swagger docs
.PHONY: swagger
swagger:
	docker run --rm -v $(PWD):/code --user $(UID) ghcr.io/swaggo/swag:v1.16.3 init --g app/api/main.go

## metrics: See metrics
.PHONY: metrics
metrics:
	expvarmon -ports="localhost:4010" \
	-vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

## grafana: Open grafana
.PHONY: grafana
grafana:
	open http://localhost:3000/

## statsviz: Open statsviz
.PHONY: statsviz
statsviz:
	open http://localhost:4010/debug/statsviz

## kafka/ui: Open kafka ui
.PHONY: kafka/ui
kafka/ui:
	open http://localhost:8080

help:
	@echo Usage:
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
