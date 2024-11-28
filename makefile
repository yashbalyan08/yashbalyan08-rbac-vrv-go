build:
	@go build -tags netgo -ldflags '-s -w' -o app

run: build
	@./docker-compose up -d
	@./app
