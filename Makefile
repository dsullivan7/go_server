DOCKER_POSTGRES = postgres:12.7
DOCKER_ALPINE = alpine:3.13.5
DOCKER_GOLANG = golang:1.16.7-alpine
DOCKER_GOLANG_LINT = golangci/golangci-lint:v1.41.1

ENVFILE ?= .env

TESTS ?= ./test/...

.PHONY: db-stop
db-stop:
	docker-compose -f docker-compose.yml stop postgres

.PHONY: db-start
db-start:
	docker-compose -f docker-compose.yml up -d postgres

.PHONY: db-create
db-create:
	docker-compose run --rm postgres-util sh -c "createdb -h \$${DB_HOST} -p \$${DB_PORT} -U \$${DB_USER} \$${DB_NAME}"

.PHONY: db-drop
db-drop:
	docker-compose run --rm postgres-util sh -c "dropdb -h \$${DB_HOST} -p \$${DB_PORT} -U \$${DB_USER} \$${DB_NAME}"

.PHONY: db-migrate
db-migrate:
	docker-compose up --build db-migrate

.PHONY: db-init
db-init:
	docker-compose up --build db-init

.PHONY: app
app-docker:
	docker-compose run --service-ports app go run ./cmd/app.go

.PHONY: app-dev
app-dev-docker:
	docker-compose up --build app-dev

.PHONY: run-docker
run-docker:
	docker-compose run --service-ports golang ./app

.PHONY: build-docker
build-docker:
	docker-compose run build go build -o bin/app ./cmd/app.go

.PHONY: test-docker
test-docker:
	docker-compose run --rm -e CGO_ENABLED=0 app go test -v $(TESTS)

.PHONY: lint-docker
lint-docker:
	docker run --rm -v ${PWD}:/data -w /data ${DOCKER_GOLANG_LINT} golangci-lint run
