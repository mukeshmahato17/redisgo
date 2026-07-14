.PHONY: redisgo

build:
	@go build -o bin/redisgo main.go

redisgo: build
	@./bin/redisgo
