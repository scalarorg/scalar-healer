ifneq (,$(wildcard .env))
    include .env
    export
endif

POSTGRES_URL := postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

migration_url := pkg/db/sqlc/migration

$(info POSTGRES_URL: $(POSTGRES_URL))
$(info migration_url: $(migration_url))


.PHONY: test coverage clean run compose down watch start stop restart shutdown deamon

test:
	ENV=test go test -coverprofile=cover.out -v ./...
coverage:
	go tool cover -html=cover.out
clean:
	rm main cover.out || true
	docker compose down --volumes --remove-orphans
down:
	docker compose down --remove-orphans
run:
	@go run cmd/api/main.go
daemon:
	@go run cmd/daemon/main.go
watch:
	air -c .air.toml
compose:
	docker compose -f compose.yml up -d --remove-orphans
start:
	docker compose -f app.compose.yml up -d --remove-orphans
stop:
	docker compose -f app.compose.yml down --remove-orphans
restart:
	docker compose -f app.compose.yml down --remove-orphans
	docker compose -f app.compose.yml up -d --remove-orphans
shutdown:
	docker compose -f app.compose.yml down --remove-orphans
	docker compose down --remove-orphans

.PHONY: sqlc migrate-up migrate-down new-migration
sqlc:
	@rm -rf pkg/db/sqlc/*.sql.go
	@./scripts/sqlc-generate.sh

migrate-up:
	migrate -path $(migration_url) -database "$(POSTGRES_URL)" -verbose up

migrate-down:
	migrate -path $(migration_url) -database "$(POSTGRES_URL)" -verbose down
new-migration:
	migrate create -ext sql -dir $(migration_url) -seq $(name)


