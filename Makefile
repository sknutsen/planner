.PHONY: help dev build sql goose-up zj live/templ live/server live/sync_assets

.DEFAULT_GOAL := help

# run templ generation in watch mode to detect all .templ files and
# re-create _templ.txt files on change, then send reload event to browser.
# Default url: http://localhost:7331
live/templ: ## Watch .templ files and reload browser via templ proxy (port 7331)
	templ generate --watch --proxy="http://127.0.0.1:8081" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
live/server: ## Rebuild and restart the app on Go changes (uses air.toml when flags omitted)
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.full_bin "CGO_ENABLED=1 BUILD_MODE=develop go build -o tmp/bin/main" \
	--build.bin "tmp/bin/main" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets: ## On assets/*.js,css changes run templ --notify-proxy for live reload
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

# start watch processes in parallel (templ, server, asset sync).
dev: ## Run live/templ, live/server, and live/sync_assets concurrently
	$(MAKE) -j3 live/templ live/server live/sync_assets

zj: ## Open zellij with layout that runs `make dev` in a pane
	zellij --layout ./zj.kdl

sql: ## Regenerate database package from sql/ (requires sqlc)
	sqlc generate

goose-up: ## Apply SQL migrations with goose (requires goose, Turso driver env)
	GOOSE_DRIVER=turso GOOSE_MIGRATION_DIR=./sql/schema/ goose up

build: ## go generate (templ) then compile the binary
	go generate ./...
	go build

help: ## Show this help
	@grep -hE '^[a-zA-Z0-9/_-]+:.*?##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
