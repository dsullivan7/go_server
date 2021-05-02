DOCKER_POSTGRES = postgres:13.2
DOCKER_ALPINE = alpine:3.13.5
DOCKER_GOLANG = golang:1.16-alpine

ENVFILE ?= .env

TESTS ?= ./test/...

.PHONY: db-run
db-run:
	docker run --name go-server-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=postgres -d $(DOCKER_POSTGRES)

.PHONY: db-remove
db-remove:
	docker rm go-server-postgres

.PHONY: db-stop
db-stop:
	docker stop go-server-postgres

.PHONY: db-start
db-start:
	docker start go-server-postgres

.PHONY: db-create
db-create:
	docker run --rm --env-file $(ENVFILE) $(DOCKER_POSTGRES) sh -c "createdb -h \$${DB_HOST} -p \$${DB_PORT} -U \$${DB_USER} \$${DB_NAME}"

.PHONY: db-drop
db-drop:
	docker run --rm --env-file $(ENVFILE) $(DOCKER_POSTGRES) sh -c "dropdb -h \$${DB_HOST} -p \$${DB_PORT} -U \$${DB_USER} \$${DB_NAME}"

.PHONY: db-migrate
db-migrate:
	docker run --rm --env-file $(ENVFILE) -v ${PWD}/internal/db/migrations:/data -w /data --entrypoint "" migrate/migrate sh -c "migrate -path /data -database postgres://\$${DB_USER}:\$${DB_PASSWORD}@\$${DB_HOST}:\$${DB_PORT}/\$${DB_NAME}?sslmode=disable up"

.PHONY: db-init
db-init:
ifeq ($(DB_DROP), yes)
	$(MAKE) db-drop
	$(MAKE) db-create
endif
	$(MAKE) db-migrate

.PHONY: app
app:
	docker-compose run --service-ports app go run ./cmd/app.go

.PHONY: deploy
deploy:
	docker run --env-file $(ENVFILE) -p 7000:7000 -v ${PWD}/bin:/data -w /data $(DOCKER_ALPINE) /data/app

.PHONY: run
run:
	docker run -t -i --env-file $(ENVFILE) -p 7000:7000 -v ${PWD}:/data -w /data $(DOCKER_GOLANG) go run ./cmd/app.go

.PHONY: build
build:
	docker run --rm --env-file $(ENVFILE) -v ${PWD}:/data -w /data $(DOCKER_GOLANG) go build -o bin/app ./cmd/app.go

.PHONY: test
test:
	docker run --rm --env-file $(ENVFILE) -v ${PWD}:/data -w /data $(DOCKER_GOLANG) go test -v $(TESTS)
