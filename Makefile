.PHONY: buildapi
buildapi: ## Build api app
	go build -o ./cmd/server/server cmd/server/main.go

.PHONY: buildagent
buildagent: ## Build agent app
	go build -o ./cmd/agent/agent cmd/agent/main.go