CONFIG_PATH := internal/config/config.json
DB_NAME := $(shell jq -r '.db_name' $(CONFIG_PATH))
DB_HOST := $(shell jq -r '.db_host' $(CONFIG_PATH))
DB_PORT := $(shell jq -r '.db_port' $(CONFIG_PATH))
DB_USER := $(shell jq -r '.db_user' $(CONFIG_PATH))
DB_PASSWORD := $(shell jq -r '.db_password' $(CONFIG_PATH))
DB_DSN := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"


.PHONY: help
help:
	@echo "Commands list:"
	@sed -n "s/^##//p" $(MAKEFILE_LIST) | column -t -s ":" | sed -e "s/^/ /"

.PHONY: db-up
db-up:
	migrate -database $(DB_DSN) -path ./migrations up

.PHONY: db-down
db-down:
	migrate -database $(DB_DSN) -path ./migrations down
