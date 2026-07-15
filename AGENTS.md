# AGENTS.md

## Cursor Cloud specific instructions

`planner` is a single Go web app (Echo + templ + HTMX) that serves a server-rendered web UI and a versioned JSON API at `/api/v1` (OpenAPI/Swagger at `/swagger`). Data lives in a Turso/libSQL (SQLite-compatible) database; auth is Auth0 (OIDC). Module: `github.com/sknutsen/planner`.

Standard commands are in the `Makefile` (`make help`): `make build`, `make test`, `make dev`, `make sql`, `make goose-up`. CI (`.github/workflows/go.yml`) only runs build + test; `go vet ./...` is the effective lint (`staticcheck` reports a few pre-existing warnings and is not wired into CI).

The update script installs `templ` (into `~/go/bin`, already on `PATH`) and runs `go mod download`. `gcc` is pre-installed; the libSQL driver needs `CGO_ENABLED=1`, so always build/test/run with CGO enabled. Toolchain is Go 1.25 (auto-downloaded via `GOTOOLCHAIN=auto`).

Note: the committed `view/*_templ.go` files carry a `templ: version: v0.3.1036` header, which is not a published templ tag (the pinned library in `go.mod` is `v0.3.1020`). Running `go generate ./...` / `make build` with an installable `templ` rewrites that header, producing harmless diffs in `view/*_templ.go`; revert them (`git checkout -- view/`) before committing unless you intend to change templates.

Non-obvious gotchas when running the app:

- **Hard startup dependencies:** `main.go` calls `os.Exit(1)` if the DB ping fails or if Auth0 OIDC discovery (`oidc.NewProvider(AUTH0_DOMAIN + "/")`) fails. Both must be reachable at boot. For real end-to-end use, supply working `AUTH0_*` and `TURSO_*` values (see `.env.example`; the app loads `.env` from the working dir, missing `.env` is only a warning).
- **Local DB without Turso Cloud:** libSQL opens a local file. Set `TURSO_DATABASE_URL=file:/absolute/path.db` with `TURSO_AUTH_TOKEN=` empty. Requires CGO.
- **Do NOT put the dev DB under `tmp/`:** `make dev` runs `air` with `clean_on_exit=true`, which wipes `tmp/` on restart (and `tmp/` is gitignored). Keep any local DB elsewhere (e.g. `/tmp/planner-local.db`).
- **Migrations are not auto-applied.** Apply `sql/schema/*.sql` (goose format, `0001`→`0009`) before use. `make goose-up` needs `goose` + a Turso driver env. If applying manually, note some files pack multiple statements inside one `-- +goose StatementBegin` block, so split on `;` and exec each (see how `database/migrate_test.go` parses goose blocks).
- **Host consistency for browser login:** when `PORT` is set but `HOST` is empty, the callback URL is built with host `0.0.0.0` while the session cookie is set on whatever host you browse (`127.0.0.1`), so the OIDC callback fails with `session state: <nil>`. Set `HOST=127.0.0.1:8081` (and `AUTH0_CALLBACK_URL=callback`) for local login.
- **Ports:** app listens on `8081` (`PORT`); `make dev` also starts a templ live-reload proxy on `7331`.
- **Fully local login without a real Auth0 tenant:** a mock OIDC provider works if its discovery `issuer` ends in a trailing slash and `AUTH0_DOMAIN` has none (`main.go` appends `/`); the `id_token` `aud` must equal `AUTH0_CLIENT_ID` and its `iss` must match the discovery `issuer`. The web session only stores the ID-token profile claims (`sub`, `name`, etc.); the access token is not re-verified for web routes (only `/api/v1` verifies Bearer JWTs against `AUTH0_API_AUDIENCE`).
