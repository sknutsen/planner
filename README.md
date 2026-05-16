# planner

## Local development

Requirements: Go 1.25+ (see `go.mod`), a C compiler for CGO (`gcc` is enough), and [templ](https://templ.guide/) on your `PATH`. Optional: [sqlc](https://sqlc.dev/) for `make sql`, [goose](https://github.com/pressly/goose) for `make goose-up`, [air](https://github.com/air-verse/air) is fetched via `go run` by the Makefile.

1. Copy `.env.example` to `.env` and set Auth0 and Turso variables.
2. Run `make help` for targets. Typical live reload: `make dev` (runs templ watch, air for Go, and air for asset changes in parallel).

The app loads `.env` from the working directory; missing `.env` only logs a warning.

### Auth0 environment variables

- **`AUTH0_AUDIENCE`**: Auth0 API identifier passed as the **`audience`** query parameter on **`/authorize`** during web login so the session access token (if used) is issued for that API. ID token verification still uses **`AUTH0_CLIENT_ID`** as required by OIDC.
- **`AUTH0_API_AUDIENCE`**: **Only** used to verify **`Authorization: Bearer`** JWTs on **`/api/v1`**. Set to the API identifier your mobile/native clients use for access tokens.

When only one API exists, these are often the **same value**; keeping both variables supports separate web vs mobile API identifiers when needed.

## Mobile JSON API (`/api/v1`)

Native and non-browser clients should call the versioned API with an **Auth0 access token** for your API (**audience** must match `AUTH0_API_AUDIENCE`). Send `Authorization: Bearer <access_token>`. The subject claim is the same `sub` used as `UserProfile.UserId` in the web app.

- **OpenAPI / Swagger UI**: [http://127.0.0.1:8081/swagger](http://127.0.0.1:8081/swagger) (adjust host/port). Raw spec: `/swagger/openapi.yaml`.
- **Sync**: list endpoints accept `updated_since`, `next_cursor`, and `limit` (default 100, max 500). Responses include `next_cursor` when another page exists.
- **Idempotency**: on mutating requests, optional `Idempotency-Key` replays the same successful JSON; a different request body with the same key returns **409** `IDEMPOTENCY_KEY_REUSE`.
- **Conflicts**: on PATCH, send `If-Match: W/"<updated_at>"` or JSON `base_updated_at`; if the row changed, the server responds with **409** and `error.entity` with the current resource.

Example:

```bash
curl -sS -H "Authorization: Bearer $ACCESS_TOKEN" http://127.0.0.1:8081/api/v1/plans
```

## Docs

- Echo - https://echo.labstack.com/docs
- Templ - https://templ.guide/
- HTMX - https://htmx.org/
- Material Icons - https://fonts.google.com/icons
- simplemde-markdown-editor -
  https://github.com/sparksuite/simplemde-markdown-editor
- markdown-it - https://github.com/markdown-it/markdown-it
- markdown-text-editor - https://github.com/nezanuha/markdown-text-editor
