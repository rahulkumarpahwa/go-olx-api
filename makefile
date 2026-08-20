build:
	@go build -o bin/api ./cmd/api
run: build
	@./bin/api
continue:
	@air
migrate-up:
	@go run ./cmd/migrate up
migrate-down:
	@go run ./cmd/migrate down