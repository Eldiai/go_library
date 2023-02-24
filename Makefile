run_postgres:
	docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d bitnami/postgresql:latest

migration_up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/golibrary?sslmode=disable" -verbose up

migration_down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/golibrary?sslmode=disable" -verbose down

run_app:
	go run cmd/api/*.go

.PHONY: run_postgres