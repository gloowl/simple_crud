psql:
	psql -h 127.0.0.1 -U admin -d simple_crud_db

build:
	go build ./src/main.go -o ./bin/crud -trimpath
run:
	go run ./src/main.go

db-up:
	docker compose up -d
db-down:
	docker compose down
