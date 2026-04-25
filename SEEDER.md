# Seeder

This project uses simple Go seed commands (no Goose integration).

## Run

1. Ensure your DB env is set (same as migrations).
2. Run migrations:

`make migrate-up`

3. Run seeders:

- Seed everything: `go run ./cmd/seed all`
- Seed menus only: `go run ./cmd/seed menu`
- Import Indonesian regions from JSON: `go run ./cmd/seed regions`

Region import reads `resources/assets` by default. Override it with `REGION_ASSETS_DIR=/path/to/assets` when needed.

## Menus

The menu seeder inserts (or updates by `slug`) the following structure:

- Dashboard
- Config
  - User
  - Role
  - Menu Role
