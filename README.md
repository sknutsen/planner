# planner

## Local development

Requirements: Go 1.25+ (see `go.mod`), a C compiler for CGO (`gcc` is enough), and [templ](https://templ.guide/) on your `PATH`. Optional: [sqlc](https://sqlc.dev/) for `make sql`, [goose](https://github.com/pressly/goose) for `make goose-up`, [air](https://github.com/air-verse/air) is fetched via `go run` by the Makefile.

1. Copy `.env.example` to `.env` and set Auth0 and Turso variables.
2. Run `make help` for targets. Typical live reload: `make dev` (runs templ watch, air for Go, and air for asset changes in parallel).

The app loads `.env` from the working directory; missing `.env` only logs a warning.

## Docs

- Echo - https://echo.labstack.com/docs
- Templ - https://templ.guide/
- HTMX - https://htmx.org/
- Material Icons - https://fonts.google.com/icons
- simplemde-markdown-editor -
  https://github.com/sparksuite/simplemde-markdown-editor
- markdown-it - https://github.com/markdown-it/markdown-it
- markdown-text-editor - https://github.com/nezanuha/markdown-text-editor
