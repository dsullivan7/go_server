docker_postgres = postgres:13.2
docker_alpine = alpine:3.13.5
docker_golang = golang:1.16-alpine

env_file = .env

ifdef ENVFILE
	env_file = $(ENVFILE)
endif

.PHONY: db-run
db-run:
	docker run --name go-server-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=postgres -d $(docker_postgres)

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
	docker run --rm --env-file $(env_file) $(docker_postgres) sh -c "createdb -h \$${DB_HOST} -p \$${DB_PORT} -U \$${DB_USER} \$${DB_NAME}"

.PHONY: db-drop
db-drop:
	docker run --rm --env-file $(env_file) $(docker_postgres) sh -c "dropdb -h \$${DB_HOST} -p \$${DB_PORT} -U \$${DB_USER} \$${DB_NAME}"

.PHONY: db-migrate
db-migrate:
	docker run --rm --env-file $(env_file) -v ${PWD}/internal/db/migrations:/data -w /data --entrypoint "" migrate/migrate sh -c "migrate -path /data -database postgres://\$${DB_USER}:\$${DB_PASSWORD}@\$${DB_HOST}:\$${DB_PORT}/\$${DB_NAME}?sslmode=disable up"

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
	docker run --env-file $(env_file) -p 7000:7000 -v ${PWD}/bin:/data -w /data $(docker_alpine) /data/app

.PHONY: run
run:
	docker run -t -i --env-file $(env_file) -p 7000:7000 -v ${PWD}:/data -w /data $(docker_golang) go run ./cmd/app.go

.PHONY: build
build:
	docker run --rm --env-file $(env_file) -v ${PWD}:/data -w /data $(docker_golang) go build -o bin/app ./cmd/app.go

.PHONY: test
test:
	docker run --rm --env-file $(env_file) -v ${PWD}:/data -w /data $(docker_golang) go test -v ./test/...
