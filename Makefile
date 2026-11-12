.PHONY: api worker worker-once frontend-install frontend-dev frontend-build build test check reset

api:
	cd backend && APP_WORKER_INLINE=true go run ./cmd/api

worker:
	cd backend && go run ./cmd/worker

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

build: frontend-build
	cd backend && go build -o bin/api ./cmd/api && go build -o bin/worker ./cmd/worker

test:
	cd backend && go test ./...

check:
	cd backend && go vet ./...
	cd frontend && npm run check

reset:
	rm -rf backend/data
