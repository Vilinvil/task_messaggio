.PHONY: lint
lint:
	golangci-lint run --timeout=3m ./...

.PHONY: easyjson
easyjson:
	easyjson pkg/myerrors/my_errors.go pkg/responses/basic_responses.go

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir ./migrations $(name)

.PHONY: mkdir-log-message
mkdir-logs:
	mkdir --parents /var/log/message

.PHONY: swag
swag:
	swag init -ot yaml --parseDependency --pdl 3 --parseInternal  --dir cmd/message/,pkg/,internal/
