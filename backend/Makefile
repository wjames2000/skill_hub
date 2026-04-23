.PHONY: build clean dev-router dev-sync dev-admin

build:
	go build -o bin/router-api ./cmd/router-api
	go build -o bin/sync-worker ./cmd/sync-worker
	go build -o bin/admin-api ./cmd/admin-api

clean:
	rm -rf bin/ tmp/

dev-router:
	go run ./cmd/router-api

dev-sync:
	go run ./cmd/sync-worker

dev-admin:
	go run ./cmd/admin-api
