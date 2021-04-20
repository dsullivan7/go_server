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

.PHONY: test
test:
	docker-compose run --rm app go test -v ./test/...
