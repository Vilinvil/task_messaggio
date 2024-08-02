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

.PHONY: test
test:
	mkdir -p bin
	 go test --race -coverpkg=./... -coverprofile=bin/cover.out ./... \
 	 && cat bin/cover.out | grep -v "mocks" | grep -v "easyjson" > bin/pure_cover.out \
  	 && go tool cover -html=bin/pure_cover.out -o=bin/cover.html \
  	 && go tool cover --func bin/pure_cover.out

.PHONY: mockgen
mockgen:
	mockgen --source=$(source) --destination=$(destination) --package=mocks

.PHONY: mockgen-all
mockgen-all:
	make mockgen source=internal/message/message/delivery/message_handler.go \
 		destination=internal/message/message/mocks/usecases.go
	make mockgen source=internal/message/message/usecases/message_service.go \
		destination=internal/message/message/mocks/reposirory.go