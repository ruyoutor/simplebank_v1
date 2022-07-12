composeup:
	docker-compose up -d

composedown:
	docker-compose down	

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine 

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	docker run --rm -v $(shell pwd)/db/migration:/migrations --network host migrate/migrate -path=/migrations/ -database postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable up

migratedown:
	docker run --rm -v $(shell pwd)/db/migration:/migrations --network host migrate/migrate -path=/migrations/ -database postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable down -all

sqlcinit:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc init

sqlc:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate	

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlcinit sqlc test server
	