# Backend System
A backend system built with Go `v1.24`.

## Introduction
Hexagonal Architecture, also known as `Ports & Adapters Architecture`, is one of several ways to build decoupled software systems. 
It was popularized by Alistair Cockburn, who is known as one of the initiators of the agile movement in software development. 
This way of organizing software is great for making applications that are easy to work on and can change without breaking.

## Project Structure
- `.docker`, contains docker files
- `.kubernetes`, contains deployment files
- `cmd`, entry points for the application
- `gen` generated rpc code
- `internal` contains the hexagonal architecture
- `pkg` core packages that are not related with the domain
- `proto` protobuf definitions
- `vendor` application dependencies

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
