DB_URL=postgresql://root:secret@localhost:5432/todo_app?sslmode=disable

postgres:
	docker run --name todo-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine

createdb:
	docker exec -it todo-postgres createdb --username=root --owner=root todo_app

dropdb:
	docker exec -it todo-postgres dropdb todo_app

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_schema:
	dbml2sql doc/db.dbml --postgres -o doc/schema.sql

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go


.PHONY: postgres createdb dropdb migrateup migratedown db_schema sqlc test server new_migration