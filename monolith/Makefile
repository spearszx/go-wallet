
build:
	@go build -o ./bin/app ./cmd/main.go

run: build
	@./bin/app

test:
	@go test -v -count=1 ./...

docker-up:
	docker-compose up -d

docker-stop:
	docker-compose down -v