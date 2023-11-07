.PHONY:

buildapi:
	go build -o ./cmd/server/server cmd/server/main.go

buildagent:
	go build -o ./cmd/agent/agent cmd/agent/main.go