
postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=Ds12345! -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

migratedown:
	migrate -path db/migration -database "postgresql://root:Ds12345!@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup:
	migrate -path db/migration -database "postgresql://root:Ds12345!@localhost:5432/simple_bank?sslmode=disable" -verbose up
sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
