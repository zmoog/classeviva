build:
	go build -o classeviva -v entrypoints/cli/main.go

mocks:
	mockery --dir adapters/spaggiari --name LoaderStorer --name Fetcher --name Adapter
	mockery --dir commands --name Command

test:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

ready: lint test build
