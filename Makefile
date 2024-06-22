postgres:
	docker run --name postgres16 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=Mwag9836 -d -p 5432:5432 postgres:16-alpine
createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres16 dropdb --username=postgres --owner=postgres simple_bank

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratestatus:
	migrate -path db/migrations -database "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable" -verbose version

migrateforce:
	migrate -path db/migrations -database "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable" -verbose force 1

sqlc:
	sqlc generate 

tests:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown migratestatus migrateforce sqlc tests server
