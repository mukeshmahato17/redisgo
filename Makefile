.PHONY: redisgo

build:
	@go build -o bin/redisgo main.go

run: build
	@./bin/redisgo
