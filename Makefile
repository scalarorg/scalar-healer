ifneq (,$(wildcard .env))
    include .env
    export
endif

POSTGRES_URL := postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

MODULE := github.com/scalarorg/scalar-healer

migration_url := pkg/db/sqlc/migration

LOCAL_LIB_PATH ?= $(shell pwd)/lib
export CGO_LDFLAGS := ${CGO_LDFLAGS} -lbitcoin_vault_ffi  -L${LOCAL_LIB_PATH}

$(info POSTGRES_URL: $(POSTGRES_URL))
$(info migration_url: $(migration_url))


.PHONY: test coverage clean run compose down watch start stop restart shutdown deamon

test:
	ENV=test go test -coverprofile=cover.out -v ./...
test-pkg:
	ENV=test go test -v $(MODULE)/pkg/$(module)/...
test-internal:
	@if [ -z "$(module)" ]; then \
		ENV=test go test -v -count=1 $(MODULE)/internal/...; \
		exit 0; \
	else \
		if [ -d "internal/$(module)" ]; then \
			ENV=test go test -v -count=1 $(MODULE)/internal/$(module)/...; \
			exit 0; \
		else \
			echo "Error: Module '$(module)' not found in internal directory"; \
			exit 1; \
		fi \
	fi
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
# @rm -rf pkg/db/sqlc/*.sql.go
sqlc:
	@./scripts/sqlc-generate.sh

migrate-up:
	migrate -path $(migration_url) -database "$(POSTGRES_URL)" -verbose up

migrate-down:
	migrate -path $(migration_url) -database "$(POSTGRES_URL)" -verbose down
new-migration:
	migrate create -ext sql -dir $(migration_url) -seq $(name)

.PHONY: lib
lib:
	mkdir -p ./lib
	cp ../bitcoin-vault/target/release/libbitcoin_vault_ffi.* ./lib

