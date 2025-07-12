include .env
export

MIGRATE=migrate
MIGRATIONS_DIR=./db/migrations
DB_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: migrate-up migrate-down migrate-force migrate-create migrate-drop

migrate-up:
	$(MIGRATE) -database "$(DB_DSN)" -path $(MIGRATIONS_DIR) up

migrate-down:
	$(MIGRATE) -database "$(DB_DSN)" -path $(MIGRATIONS_DIR) down

migrate-force:
	$(MIGRATE) -database "$(DB_DSN)" -path $(MIGRATIONS_DIR) force

migrate-drop:
	$(MIGRATE) -database "$(DB_DSN)" -path $(MIGRATIONS_DIR) drop

migrate-create:
ifndef name
	$(error "❌ You must pass a name: make migrate-create name=create_articles_table")
endif
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

.PHONY: seed unit-test
seed:
	go run ./cmd/seed/main.go

unit-test:
	go test ./... -v -coverprofile=coverage.out

