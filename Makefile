.DEFAULT_GOAL := build

GM_HOST = "localhost"
GM_PORT = "8080"

PG_USER = "gophermart-user"
PG_PASSWORD = "gophermart-password"
PG_DB = "db-gophermart"
PG_HOST = "localhost"
PG_PORT = "5432"
PG_DATABASE_DSN = "postgres://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB}?sslmode=disable"
PG_IMAGE = "postgres:13.13-bullseye"
PG_DOCKER_CONTEINER_NAME = "gophermart-pg-13.3"

.PHONY:build
build:
	go build -C ./cmd/gophermart/ -o gophermart

.PHONY:clean
clean:
	-rm -f ./cmd/gophermart/gophermart

.PHONY:statictest
statictest:
	go vet -vettool=$$(which statictest) ./...

.PHONY:test
test: build statictest
	go test -v ./...

.PHONY:gophermart-run
gophermart-run: build 
	chmod +x ./cmd/gophermart/gophermart && \
	./cmd/gophermart/gophermart -a ${GM_HOST}:${GM_PORT} -d ${PG_DATABASE_DSN}

.PHONY:gophermarttest
gophermarttest: test db-up
	./cmd/gophermart/gophermarttest \
		-test.v -test.run=^TestGophermart$ \
		-gophermart-binary-path=cmd/gophermart/gophermart \
		-gophermart-host=${GM_HOST} \
		-gophermart-port=${GM_PORT} \
		-gophermart-database-uri=${PG_DATABASE_DSN} \
		-accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
		-accrual-host=localhost \
		-accrual-port=$$(random unused-port) \
		-accrual-database-uri=${PG_DATABASE_DSN}

.PHONY:gophermarttest-user-auth
gophermarttest-user-auth: test db-up
	./cmd/gophermart/gophermarttest \
		-test.v -test.run=^TestGophermart/TestUserAuth \
		-gophermart-binary-path=cmd/gophermart/gophermart \
		-gophermart-host=${GM_HOST} \
		-gophermart-port=${GM_PORT} \
		-gophermart-database-uri=${PG_DATABASE_DSN} \
		-accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
		-accrual-host=localhost \
		-accrual-port=$$(random unused-port) \
		-accrual-database-uri=${PG_DATABASE_DSN}

.PHONY: db-up
db-up:
	PG_USER=${PG_USER} \
	PG_PASSWORD=${PG_PASSWORD} \
	PG_DB=${PG_DB} \
	PG_HOST=${PG_HOST} \
	PG_PORT=${PG_PORT} \
	PG_DATABASE_DSN=${PG_DATABASE_DSN} \
	PG_IMAGE=${PG_IMAGE} \
	PG_DOCKER_CONTEINER_NAME=${PG_DOCKER_CONTEINER_NAME} \
	docker compose -f ./docker-compose.yml up -d postgres

.PHONY: db-down
db-down:
	PG_USER=${PG_USER} \
	PG_PASSWORD=${PG_PASSWORD} \
	PG_DB=${PG_DB} \
	PG_HOST=${PG_HOST} \
	PG_PORT=${PG_PORT} \
	PG_DATABASE_DSN=${PG_DATABASE_DSN} \
	PG_IMAGE=${PG_IMAGE} \
	PG_DOCKER_CONTEINER_NAME=${PG_DOCKER_CONTEINER_NAME} \
	docker compose -f ./docker-compose.yml down postgres

GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

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

.PHONY: golangci-lint-clean
golangci-lint-clean:
	rm -rf ./golangci-lint

