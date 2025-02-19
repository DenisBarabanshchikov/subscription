run:
	godotenv -f .env go run ./cmd/main.go
doc:
	swag init -g cmd/main.go