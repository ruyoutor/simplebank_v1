postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine 

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	docker run -v $(shell pwd)/db/migration:/migrations --network host migrate/migrate -path=/migrations/ -database postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable up
	# migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
	# docker run -v $(shell pwd)/db/migration:/migrations --network host migrate/migrate -path=/migrations/ -database postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable down 1

sqlcinit:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc init

sqlc:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate	

.PHONY: postgres createdb dropdb migrateup migratedown, sqlcinit, sqlc
	