# Docker commands
create_postgres:
	docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=password -d postgres:13-alpine

drop_postgress:
	docker rm postgres13

# Postgres command
create_db:
	docker exec -it postgres13 createdb --username=admin --owner=admin freefortalking 

drop_db:
	docker exec -it postgres13 dropdb -U admin freefortalking

# Database migrations
migration_up:
	migrate -path db/migration -database "postgres://admin:password@localhost:5432/freefortalking?sslmode=disable" -verbose up

migration_down:
	migrate -path db/migration -database "postgres://admin:password@localhost:5432/freefortalking?sslmode=disable" -verbose down

.PHONY: create_postgres drop_postgress create_db drop_db migration_up migration_down