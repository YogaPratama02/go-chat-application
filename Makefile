Arguments := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

setup:
	go mod tidy

env:
	cp .env.example .env

server:
	go run main.go

create-table:
	migrate create -ext sql -dir migrations -seq create_$(Arguments)_table

migrate-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/chat_application?sslmode=disable" -path migrations up