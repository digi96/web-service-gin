migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/contact_db?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/contact_db?sslmode=disable" --verbose down

sqlc:
	sqlc generate

swag:
	swag init

test:
	go test -v -cover ./...

rabbitmq:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management


.PHONY: migrateup migratedown sqlc swag test rabbitmq