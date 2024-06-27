postgres:
	docker run --name postgres16 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=Mwag9836 -d -p 5432:5432 postgres:16-alpine
createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres16 dropdb --username=postgres --owner=postgres simple_bank

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migrations -database "postgresql://postgres:Mwag9836@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

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

mock:
	mockgen -build_flags=--mod=mod -package mockdb  -destination db/mock/store.go github.com/joshuandeleva/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migratestatus migrateforce sqlc tests server mock migrateup1 migratedown1
