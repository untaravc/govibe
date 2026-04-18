# Database migrations (MySQL)

This repo uses `goose` for migrations, with migration files stored in `database/migrations/`.

## Configure

Set DB config in `.env` (see `.env.example`).

## Create a migration

```bash
go run ./cmd/migrate create create_users_table
```

This creates a timestamped SQL file in `database/migrations/` with `-- +goose Up` / `-- +goose Down` sections.

## Run migrations

Apply all pending migrations:

```bash
go run ./cmd/migrate up
```

Rollback the last migration:

```bash
go run ./cmd/migrate down
```

Show status:

```bash
go run ./cmd/migrate status
```

## Versions

Migrations are versioned by the timestamp prefix in the filename (e.g. `20260418123456_create_users_table.sql`).

Run up/down to a specific version:

```bash
go run ./cmd/migrate up-to 20260418123456
go run ./cmd/migrate down-to 20260418123456
```

