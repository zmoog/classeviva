build:
	go build -o classeviva -v entrypoints/cli/main.go

lint:
	golangci-lint run ./...