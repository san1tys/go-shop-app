.PHONY: dev backend frontend up down

dev:
	cd deployments && docker compose up --build

up:
	cd deployments && docker compose up -d --build

down:
	cd deployments && docker compose down

backend-local:
	cd backend && go run ./cmd/api

frontend-local:
	cd frontend && npm run dev
