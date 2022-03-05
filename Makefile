include .env

.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-new test test-ci test-spec start-db

migrate-up:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} up

migrate-down:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} down ${n}

migrate-force:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} force ${v}

migrate-version:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} version

migrate-new:
	@migrate create -ext sql -dir ${MIGRATION_PATH} -seq ${name}

test:
	@go test -failfast -parallel 2 ./...

test-ci:
	@go test -failfast -cover -v -parallel 2 ./...

test-spec:
	APIURL=${API_URL} bash ./conduit/spec/run-api-tests.sh

start-db:
	@docker-compose --env-file .env up -d