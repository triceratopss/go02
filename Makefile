include .env

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: db_migrate
db_migrate:
	docker compose run --rm go02-db-migration up
