# TODO App

A small task manager with a Go + Postgres backend and a dependency-free HTML/JS frontend.

**Live:** https://todo-app-vi2s.onrender.com

> Hosted on Render's free tier, which sleeps after ~15 minutes of inactivity. The first request after a nap takes ~30 seconds while the service wakes up. Subsequent requests are instant.

## Architecture

```mermaid
flowchart LR
    B[Browser<br/>plain HTML + JS] -->|fetch, same origin| G[Gin router]
    G --> H[task handlers<br/>validate, map errors]
    H --> Q[sqlc queries]
    Q --> PG[(Postgres<br/>tasks table)]
```

A task's life: created → (completed ⇄ not completed) → soft-deleted (`deleted_at` set, hidden from every read).

## Stack

| Layer | Choice | Why |
|---|---|---|
| Language | Go 1.26 | Type safety, single static binary, fast startup |
| HTTP | [Gin](https://github.com/gin-gonic/gin) | JSON binding and request validation out of the box |
| Database | PostgreSQL 18 | Data survives restarts; real constraints and defaults |
| Driver | [pgx/v5](https://github.com/jackc/pgx) | Modern Postgres driver; `pgtype` models nullable columns correctly |
| Queries | [sqlc](https://sqlc.dev) | Hand-written SQL compiled into type-safe Go; query errors surface at build time, not runtime |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) | Numbered, ordered up/down migrations |
| Frontend | Plain HTML, CSS, JS | No framework. `fetch()` against the same origin, so no CORS |

## API

All responses are JSON. Errors are `{"error": "..."}`.

| Method | Path | Body | Success | Errors |
|---|---|---|---|---|
| `POST` | `/tasks` | `{"title": "..."}` | `201` + task | `400` empty or >255 chars |
| `GET` | `/tasks` | — | `200` + `[task]` | — |
| `GET` | `/tasks/{id}` | — | `200` + task | `400` bad id, `404` |
| `PATCH` | `/tasks/{id}` | `{"title": "..."}` | `200` + task | `400`, `404` |
| `PATCH` | `/tasks/{id}/complete` | `{"completed": bool}` | `200` + task | `400`, `404` |
| `DELETE` | `/tasks/{id}` | — | `204` | `400` bad id |

A task looks like:

```json
{
  "id": 1,
  "title": "buy milk",
  "completed": false,
  "created_at": "2026-08-27T13:55:36.917611Z",
  "updated_at": "2026-08-27T13:56:13.42521Z",
  "deleted_at": null
}
```

## Running locally

Requires Go, Docker, [`migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate), and [`sqlc`](https://docs.sqlc.dev/en/latest/overview/install.html).

```bash
make postgres     # start Postgres in Docker (container: todo-postgres)
make createdb     # create the todo_app database
make migrateup    # apply migrations
make server       # run on :8080
```

Open http://localhost:8080.

Configuration is read from the environment, with local defaults:

| Variable | Default |
|---|---|
| `DATABASE_URL` | `postgresql://root:secret@localhost:5432/todo_app?sslmode=disable` |
| `PORT` | `8080` |

Other targets: `make sqlc` (regenerate query code), `make new_migration name=<name>`, `make migratedown`, `make dropdb`.

## Design decisions

**Soft delete via a nullable timestamp.** `DELETE /tasks/{id}` never removes a row. It sets `deleted_at = now()`. Every read query carries `WHERE deleted_at IS NULL`, so deleted tasks are invisible to the API but remain in the database for undo and audit. A nullable timestamp was chosen over an `is_deleted` boolean because it answers both *whether* and *when* with one column. The cost is discipline: forgetting the filter in any query resurrects deleted rows.

**Completed is separate from deleted.** `completed` is a user-facing state (show it with a strikethrough). `deleted_at` is a system state (hide it). They are independent columns because they mean different things.

**URLs are nouns; methods are verbs.** `POST /tasks`, not `/create_task`. The HTTP method already says what is happening.

**The server owns its metadata.** Clients send only `title` (and `completed` for the toggle). `id`, `created_at`, `updated_at`, and `deleted_at` are set by the database or the server. A client cannot forge a creation time or pick its own id.

**`PATCH`, not `PUT`.** `PUT` semantically replaces the whole resource, so omitted fields get wiped. `PATCH` changes only what is sent.

**Mutations return the row.** `POST` and both `PATCH` endpoints respond with the task as the database now sees it, so the client never has to guess at server-set fields.

**Soft delete is an implementation detail.** The API exposes a normal `DELETE` and returns `204`. Callers do not need to know the row still exists.

**Errors are translated, not leaked.** Database errors are logged server-side and returned as a generic `500`. `pgx.ErrNoRows` becomes `404`. A malformed id is `400`. The raw error never reaches the client.

**Config from the environment.** `DATABASE_URL` and `PORT` are read from env vars with local fallbacks, so the same binary runs locally and in production.

## Schema

```sql
CREATE TABLE tasks (
  id          bigserial PRIMARY KEY,
  title       varchar(255) NOT NULL,
  completed   bool         NOT NULL DEFAULT false,
  created_at  timestamptz  NOT NULL DEFAULT now(),
  updated_at  timestamptz  NOT NULL DEFAULT now(),
  deleted_at  timestamptz                          -- NULL means alive
);
```

Source of truth is [`db/migration/`](db/migration/). The DBML in [`doc/`](doc/) is a visual companion.

## Layout

```
.
├── main.go            config, DB pool, wires the server
├── api/
│   ├── server.go      Server struct, routes
│   └── task.go        handlers
├── db/
│   ├── migration/     golang-migrate up/down files (append-only)
│   ├── query/         hand-written SQL, input to sqlc
│   └── sqlc/          generated Go (do not edit)
├── doc/               DBML schema + generated SQL snapshot
├── static/            index.html, style.css
├── Makefile
└── sqlc.yaml
```

## Not in v1

Deliberately left out to ship on time. Roughly in the order I'd add them:

- A `TaskStore` interface so handlers don't depend on `*db.Queries`, plus an in-memory implementation for tests
- Handler tests
- Pagination and filtering on `GET /tasks`
- Partial index on `tasks (id) WHERE deleted_at IS NULL`
- Fail on missing `DATABASE_URL` in production instead of falling back to localhost
- Structured logging, Dockerfile, auth
