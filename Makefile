POSTGRES_URL="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
MIGRATION_PATH="./db/migrations"
API_URL="http://localhost:4000/api"

.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-new test test-ci test-spec

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
	@go test -parallel 2 ./...

test-ci:
	@go test -cover -v -parallel 2 ./...

test-spec:
	APIURL=${API_URL} bash ./conduit/spec/run-api-tests.sh