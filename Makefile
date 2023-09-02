# Variables
FRONTEND_DIR=frontend
BACKEND_BIN=bin/trophies
GO_MAIN=main.go
GO_FLAGS=-ldflags '-w -s'
NODE_ENV=production

# Build frontend
build-frontend:
	NODE_ENV=$(NODE_ENV) pnpm --dir $(FRONTEND_DIR) build

# Build backend
build-backend:
	CGO_ENABLED=0 go build -o $(BACKEND_BIN) $(GO_FLAGS) $(GO_MAIN)

# Run migrations
migrate:
	go run $(GO_MAIN) migrate

# Run application
run: build-frontend migrate
	go run $(GO_MAIN) serve

# Fetch data (example target)
fetch: migrate
	go run $(GO_MAIN) fetch

# Run in production mode
run-prod: build-frontend build-backend
	./$(BACKEND_BIN) serve
