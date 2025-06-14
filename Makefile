# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

GO_VERSION := 1.24
UID := $(shell id -u)
GID := $(shell id -g)

INPUT ?= $(shell bash -c 'read -p "Insert name: " name; echo $$name')
CURRENT_TIME := $(shell date --iso-8601=seconds)
GIT_VERSION := $(shell git describe --always --dirty --tags --long)
LINKER_FLAGS := "-s -X main.buildTime=${CURRENT_TIME} -X main.version=${GIT_VERSION}"

DOCKER_COMPOSE_LOCAL := docker compose -f ./compose.yml
MIGRATE := $(DOCKER_COMPOSE_LOCAL) run --rm migrate
MIGRATION_DB_DSN := "postgres://housi:secret123@db:5432/housi_db?sslmode=disable"

.PHONY: help
help:
	@echo Usage:
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## ==================
## Docker
## ==================

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
	make docker/down
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
.PHONY: go/http/run
go/http/run:
	go run cmd/http/main.go

## go/http/build: build the http application
.PHONY: go/http/build
go/http/build:
	cd cm/http & \
	go build -ldflags=${LINKER_FLAGS} -o=./http-api

## ==================
## gRPC Application
## ==================

## go/grpc/run: Run gRPC locally
.PHONY: go/grpc/run
go/grpc/run:
	go run cmd/grpc/main.go

## go/grpc/build: build the gRPC application
.PHONY: go/grpc/build
go/grpc/build:
	cd cmd/grpc & \
	go build -ldflags=${LINKER_FLAGS} -o=./banking-api

## go/grpc/curl: curl the gRPC application
.PHONY: go/grpc/curl
go/grpc/curl:
	grpcurl -plaintext -format json -d '{"id": "45b5fb22-755f-aa79-8a07-a58d42a1fa2f"}' localhost:9090 user.v1.UserService/GetUserById

## go/grpc/curl/list: List all endpoints for gRPC application
.PHONY: go/grpc/curl/list
go/grpc/curl/list:
	grpcurl -plaintext localhost:9090 list

## ==================
## CLI Application
## ==================

## go/cli/build: Build cli application
.PHONY: go/cli/build
go/cli/build:
	go build -o cmd/cli/cli cmd/cli/main.go

## go/cli/genkey: Generate key
.PHONY: go/cli/genkey
go/cli/genkey:
	make go/cmd/build
	app/cmd/cmd genkey

## go/cli/useradd: Add user
.PHONY: go/cli/useradd
go/cli/user/add:
	make go/cli/build
	cmd/cli/cli useradd "chris housi" "example@example.com" "1232455477"

## go/cli/user/events: User events
.PHONY: go/cli/userevents
go/cli/user/events:
	make go/cli/build
	cmd/cli/cli userevents

## ==================
## Database
## ==================

## repository/migrations/create name=$1: Create new migration files
.PHONY: db/migrate/create
db/migrate/create:
	$(MIGRATE) create -seq -ext=.sql -dir=./database/migrations $(INPUT)

## repository/migrations/up: Apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up:
	$(MIGRATE) -path=./database/migrations -database=${MIGRATION_DB_DSN} up

## repository/migrations/down: Apply all down database migrations (DROP Database)
.PHONY: db/migrate/down
db/migrate/down:
	$(MIGRATE) -path=./database/migrations -database=${MIGRATION_DB_DSN} down

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

## errcheck: Check for unhandled errors
.PHONY: errcheck
errcheck:
	go tool errcheck ./...

## test: Run tests
.PHONY: test
test:
	CGO_ENABLED=1 go test -v -cover -short -race -json -p 4 ./... | go tool tparse --all

## coverage: Inspect coverage
.PHONY: coverage
coverage:
	go test -cover -v -coverpkg=./... ./...
	go test -coverprofile=coverage.out -coverpkg=./... ./...
	grep -Ev "test/|gen/|debug/|dbtest|unitest" coverage.out > filtered.out
	go tool cover -func=filtered.out

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
## Buf
## ==================

## buf/init: Buf initialize configuration
.PHONY: buf/init
buf/init:
	go tool buf config init

## buf/lint: Buf linter
.PHONY: buf/lint
buf/lint:
	go tool buf lint .

## buf/format: Format protobuf
.PHONY: buf/generate
buf/format:
	go tool buf format .

## buf/generate: Generate protobuf
.PHONY: buf/generate
buf/generate:
	go tool buf generate .

## ==================
## Tooling
## ==================

tools/install:
	go install tool

tools/list:
	go tool

tools/update:
	go get -u tool

## ==================
## Utils
## ==================

# swagger: Generate swagger docs
.PHONY: swagger
swagger:
	docker run --rm -v $(PWD):/code ghcr.io/swaggo/swag:v1.16.3 init --g app/http/main.go

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
