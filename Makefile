.PHONY: lint
lint:
	golangci-lint run --fix --timeout=3m ./...

.PHONY: easyjson
easyjson:
	easyjson pkg/myerrors/my_errors.go pkg/responses/basic_responses.go pkg/models/counter_message.go

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir ./migrations $(name)

.PHONY: mkdir-log-message
mkdir-logs:
	mkdir --parents /var/log/message

.PHONY: mkdir-log-message_worker
mkdir-log-message_worker:
	 mkdir --parents /var/log/message_worker

.PHONY: swag
swag:
	swag init -ot go,yaml --parseDependency --parseInternal  --dir cmd/message/,pkg/,internal/
