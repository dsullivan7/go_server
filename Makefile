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
	docker run --rm --env-file .env --network="host" postgres:13.2-alpine sh -c "createdb -h \$${DB_HOST} -p \$${DB_PORT} -U \$${DB_USER} \$${DB_NAME}"

.PHONY: db-drop
db-drop:
	docker run --rm --env-file .env --network="host" postgres:13.2-alpine sh -c "dropdb -h \$${DB_HOST} -p \$${DB_PORT} -U \$${DB_USER} \$${DB_NAME}"

.PHONY: db-migrate
db-migrate:
	docker run --rm --env-file .env -v ${PWD}/internal/db/migrations:/data -w /data --network="host" --entrypoint "" migrate/migrate sh -c "migrate -path /data -database postgres://\$${DB_USER}:\$${DB_PASSWORD}@\$${DB_HOST}:\$${DB_PORT}/\$${DB_NAME}?sslmode=disable up"

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
