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

.PHONY: runapi
runapi: ## Run api app
	go run -ldflags  "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=v1" cmd/server/main.go

.PHONY: buildagent
buildagent: ## Build agent app
	go build -o ./cmd/agent/agent cmd/agent/main.go

.PHONY: buildcustomlinter
buildcustomlinter: ## Build linter
	go build -o ./cmd/staticlint/staticlint cmd/staticlint/main.go

.PHONY: rundb
rundb:
	docker run --name metrics -e POSTGRES_USER=usr -e POSTGRES_PASSWORD=psswrd -e POSTGRES_DB=metrics -p 54323:5432 -d postgres:14.2

.PHONY: migrate
migrate:
	goose -dir=internal/repo/internal/db/migrations postgres "host=localhost port=54323 user=usr password=psswrd dbname=metrics sslmode=disable" up

.PHONY: rollbackmigrations
rollbackmigrations:
	goose -dir=internal/repo/internal/db/migrations postgres "host=localhost port=54323 user=usr password=psswrd dbname=metrics sslmode=disable" reset


## TESTS

MOCKS_DESTINATION=mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: ./internal/app/agent/imetric.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done

.PHONY: test
test:
	go install gotest.tools/gotestsum@latest
	gotestsum --format pkgname -- -coverprofile=cover.out ./...

.PHONY: bench
bench:
	go test -bench . -benchmem ./...

.PHONY: cover
cover:
	go test -v -coverpkg=./... -coverprofile=cover.out.tmp ./...
	cat cover.out.tmp | grep -v "_easyjson.go" | grep -v "/mocks/" | grep -v "/db/"> cover.out
	rm cover.out.tmp
	go tool cover -func cover.out

## LINTERS
.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w -local github.com/NStegura/metrics ./cmd
	goimports -w -local github.com/NStegura/metrics ./internal

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml --out-format=colored-line-number --sort-results
	./cmd/staticlint/staticlint ./...

## PROFILE

.PHONY: pprofcpu
pprofcpu:
	go tool pprof -http=":9090" -seconds=30 http://localhost:8081/debug/pprof/profile

.PHONY: pprofmem
pprofmem:
	go tool pprof -http=":9090" -seconds=30 http://localhost:8081/debug/pprof/heap

.PHONY: pprofmemfile # save to file and check
pprofmemfile:
	curl -sK -v http://localhost:8081/debug/pprof/heap > heap.out
	go tool pprof -http=":9090" -seconds=30 http://localhost:8081/debug/pprof/heap

.PHONY: pprofconsolecpu # save to file and check
pprofconsolecpu:
	go tool pprof -seconds=30 http://localhost:8081/debug/pprof/profile

.PHONY: pprofsavemem # example save to file
pprofsavemem:
	curl http://localhost:8081/debug/pprof/heap > ./profiles/result.pprof

.PHONY: pprofcompare # example compare
pprofcompare:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof