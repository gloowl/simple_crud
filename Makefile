build:
	go build -trimpath -o ./bin/herbs-cli ./src 
run:
	go run ./src/main.go

db-up:
	docker compose up -d
db-down:
	docker compose down
