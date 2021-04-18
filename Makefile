db-create:
	docker-compose run --rm db createdb

db-drop:
	docker-compose run --rm db dropdb go_server

db-migrate:
	docker-compose run --rm db-migrate

db-init:
ifeq ($(DB_DROP), yes)
	$(MAKE) db-drop
	$(MAKE) db-create
endif
	$(MAKE) db-migrate
