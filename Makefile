build:
	go build ./src/main.go -o ./bin/crud -trimpath
run:
	go run ./src/main.go

db-up:
	docker compose up -d
db-down:
	docker compose down
