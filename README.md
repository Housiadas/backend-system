# Backend System
A backend system built with Go `v1.24`.

## Introduction
The system is separated in three main layers
- `app` layer, represents the entry point of our application (request/response)
- `business` layer, represents the business logic, like communication with database
- `foundation` layer, represents core modules of the application, like logger

## App Layer
- `http`, REST API implementation
- `gRPC` implementation
- `cmd` commands implementation
- `domain` contains the entry point for our business logic

## Docker

```
make docker/up
make docker/down
```

## Migrations

```
make db/migrate/up
make db/migrate/down
make go/cmd/seed
```

## Run tests

```
make tests
```

## Help

```
make help
```

## Database Schema
We used https://dbdiagram.io to define our database structure.
