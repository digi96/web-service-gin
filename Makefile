migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/contact_db?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/contact_db?sslmode=disable" --verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...


.PHONY: migrateup migratedown sqlc test