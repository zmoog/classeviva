build:
	go build -o classeviva -v entrypoints/cli/main.go

test:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

ready: lint test build
