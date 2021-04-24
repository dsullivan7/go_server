.PHONY: db-run
db-run:
	docker run --name go-server-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=$$(whoami) -d postgres:13.2-alpine

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
	docker-compose run --rm db createdb

.PHONY: db-drop
db-drop:
	docker-compose run --rm db dropdb go_server

.PHONY: db-migrate
db-migrate:
	docker-compose run --rm db-migrate

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

.PHONY: build
build:
	docker run --rm --env-file .env -v ${PWD}:/data -w /data golang:1.16-alpine go build -o bin/app ./cmd/app.go

.PHONY: test
test:
	docker run --rm --env-file .env -v ${PWD}:/data -w /data golang:1.16-alpine go test -v ./test/...
