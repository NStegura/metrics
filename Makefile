$(eval CNT_API := $(shell docker ps -f name=metrics-api -q | wc -l | awk '{print $1}'))

.PHONY: up
up:
ifeq ($(CNT_API),0)
	docker-compose up --build --no-recreate --detach; sleep 5
endif

.PHONY: bash
bash:
	make up
	docker-compose exec metrics-api /bin/sh

.PHONY: down
down:
	docker-compose down --remove-orphans --rmi local

.PHONY: buildall
buildall: buildapi buildagent

.PHONY: buildapi
buildapi: ## Build api app
	go build -o ./cmd/server/server cmd/server/main.go

.PHONY: buildagent
buildagent: ## Build agent app
	go build -o ./cmd/agent/agent cmd/agent/main.go

.PHONY: rundb
rundb:
	docker run --name metrics -e POSTGRES_USER=usr -e POSTGRES_PASSWORD=psswrd -e POSTGRES_DB=metrics -p 54323:5432 -d postgres:14.2

## LINTERS
GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: lint
lint: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.55.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint
