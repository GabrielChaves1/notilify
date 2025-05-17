include ./cmd/api/.env
export

migrate-create:
	migrate create -ext sql -dir ./migrations -seq ${name}
migrate-up:
	migrate -path ./migrations -database $(POSTGRES_CONNECTION_STRING) up
migrate-down:
	migrate -path ./migrations -database $(POSTGRES_CONNECTION_STRING) down
