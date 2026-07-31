include .envrc
DB_ADDR ?= postgres://postgres:postgres@localhost/social?sslmode=disable

MIGRATION_PATH := ./cmd/migrate/migrations

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_ADDR)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_ADDR)" down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-version
migrate-version:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_ADDR)" version

.PHONY: migrate-drop
migrate-drop:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_ADDR)" drop -f

.PHONY: migrate-force
migrate-force:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_ADDR)" force $(version)
	
%:
	@: