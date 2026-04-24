.PHONY: build clean dev-router dev-sync dev-admin

CONFIG ?= config.yaml

build:
	go build -o bin/router-api ./cmd/router-api
	go build -o bin/sync-worker ./cmd/sync-worker
	go build -o bin/admin-api ./cmd/admin-api

clean:
	rm -rf bin/ tmp/

dev-router:
	go run ./cmd/router-api --config=$(CONFIG)

dev-sync:
	go run ./cmd/sync-worker --config=$(CONFIG)

dev-admin:
	go run ./cmd/admin-api --config=$(CONFIG)
