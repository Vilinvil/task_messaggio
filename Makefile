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
mkdir-log-message:
	mkdir --parents /var/log/message

.PHONY: mkdir-log-messageworker
mkdir-log-messageworker:
	 mkdir --parents /var/log/messageworker

.PHONY: swag
swag:
	swag init -ot go,yaml --parseDependency --parseInternal  --dir cmd/message/,pkg/,internal/

.PHONY: cp-env
cp-env:
	cp -r .env.example/ .env/
