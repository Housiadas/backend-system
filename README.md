# Backend System
A backend system built with Go `v1.23`.

## Setup

```
make docker/build
make docker/up
make db/migrate/up
make go/cmd/seed
```

## Migration

```
make db/migrate/create
make db/migrate/up
make db/migrate/down
```

## Testing

```
make tests
```

## Help

```
make help
```

## Database Schema
We used https://dbdiagram.io to define our database structure.
### Needs update!
