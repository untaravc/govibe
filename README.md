# govibe

Fullstack starter: Go + Fiber, server-rendered templates, and static assets.

## Run

```bash
go mod tidy
go run ./cmd/server
```

Or:

```bash
make run
```

## Dev (Air)

```bash
air -c air.toml
```

## Frontend (Vue)

Source files live in `resources/js/`. `npm run dev` compiles Vue into `public/dist/`, and `resources/views/layout.html` loads `/static/dist/app.js`.

```bash
npm install
npm run dev
```

## Tailwind CSS

Tailwind is compiled by Vite into `public/dist/app.css` and loaded by `resources/views/layout.html`.

Open `http://localhost:3000`.

## Database config

MySQL defaults are loaded from `.env` (see `.env.example`). Config loader: `configs.LoadDatabaseConfig()`.
