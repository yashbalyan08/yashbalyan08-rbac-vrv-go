build:
	@docker-compose up -d
	@go build -tags netgo -ldflags '-s -w' -o app

run: build
	@./app
