.PHONY: lint
lint:
	golangci-lint run --timeout=3m ./...

.PHONY: easyjson
easyjson:
	easyjson pkg/myerrors/my_errors.go pkg/responses/basic_responses.go
