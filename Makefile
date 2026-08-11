-include .env

DB_DSN := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE)"


.PHONY: help
help:
	@echo "Commands list:"
	@sed -n "s/^##//p" $(MAKEFILE_LIST) | column -t -s ":" | sed -e "s/^/ /"


## db-up: Runs database migrations up.
.PHONY: db-up
db-up:
	migrate -database $(DB_DSN) -path ./migrations up

## db-down: Rolls back database migrations down.
.PHONY: db-down
db-down:
	migrate -database $(DB_DSN) -path ./migrations down


## server-up: Creates and launches docker containers with server and database.
.PHONY: server-up
server-up:
	docker compose up -d

## server-down: Shuts docker containers with server and database.
.PHONY: server-down
server-down:
	docker compose down
