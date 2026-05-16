# planner

## Local development

Requirements: Go 1.25+ (see `go.mod`), a C compiler for CGO (`gcc` is enough), and [templ](https://templ.guide/) on your `PATH`. Optional: [sqlc](https://sqlc.dev/) for `make sql`, [goose](https://github.com/pressly/goose) for `make goose-up`, [air](https://github.com/air-verse/air) is fetched via `go run` by the Makefile.

1. Copy `.env.example` to `.env` and set Auth0 and Turso variables.
2. Run `make help` for targets. Typical live reload: `make dev` (runs templ watch, air for Go, and air for asset changes in parallel).

The app loads `.env` from the working directory; missing `.env` only logs a warning.

## Mobile JSON API (`/api/v1`)

Native and non-browser clients should call the versioned API with an **Auth0 access token** for your API (**audience** must match `AUTH0_AUDIENCE`). Send `Authorization: Bearer <access_token>`. The subject claim is the same `sub` used as `UserProfile.UserId` in the web app.

- **Sync**: list endpoints accept `updated_since`, `cursor_ts`, `cursor_id`, and `limit` (default 200, max 500). Responses may include `next_cursor` with `ts` and `id` for the next page.
- **Idempotency**: for `POST`/`PATCH`/`DELETE` on tasks and plans, send `Idempotency-Key` (any unique string per logical request); retries return the same JSON body and status.
- **Conflicts**: for `PATCH /plans/:planId` and `PATCH /tasks/:id`, include `base_updated_at` from the last read copy of the resource. If the row changed, the server responds with **409** and `error.current` holding the latest JSON.

Example:

```bash
curl -sS -H "Authorization: Bearer $ACCESS_TOKEN" http://127.0.0.1:8081/api/v1/me
```


- Echo - https://echo.labstack.com/docs
- Templ - https://templ.guide/
- HTMX - https://htmx.org/
- Material Icons - https://fonts.google.com/icons
- simplemde-markdown-editor -
  https://github.com/sparksuite/simplemde-markdown-editor
- markdown-it - https://github.com/markdown-it/markdown-it
- markdown-text-editor - https://github.com/nezanuha/markdown-text-editor
