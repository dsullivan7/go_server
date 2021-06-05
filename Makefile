DOCKER_POSTGRES = postgres:13.2
DOCKER_ALPINE = alpine:3.13.5
DOCKER_GOLANG = golang:1.16-alpine

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
	docker-compose run --rm migrate

.PHONY: db-init
db-init:
ifeq ($(DB_DROP), yes)
	$(MAKE) db-drop
	$(MAKE) db-create
endif
	$(MAKE) db-migrate

.PHONY: run
run:
	docker-compose run --service-ports app go run ./cmd/app.go

.PHONY: deploy
deploy:
	docker-compose run --service-ports deploy ./app

.PHONY: build
build:
	docker-compose run build go build -o bin/app ./cmd/app.go

.PHONY: test
test:
	$(MAKE) db-drop
	$(MAKE) db-create
	$(MAKE) db-migrate
	docker-compose run --rm app go test -v $(TESTS)
