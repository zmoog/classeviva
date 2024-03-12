build:
	go build -o classeviva -v entrypoints/cli/main.go

mock:
	mockery --dir adapters/spaggiari --name Client
	mockery --dir adapters/spaggiari --name Fetcher
	mockery --dir adapters/spaggiari --name GradesReceiver
	mockery --dir adapters/spaggiari --name HTTPClient
	mockery --dir adapters/spaggiari --name LoaderStorer
	mockery --dir adapters/spaggiari --name NoticeboardsReceiver

	mockery --dir commands --name Command

test:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

ready: lint test build
