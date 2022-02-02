POSTGRES_URL="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
MIGRATION_PATH="./db/migrations"
API_URL="http://localhost:4000/api"

migrate-up:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} up

migrate-down:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} down ${n}

migrate-force:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} force ${v}

migrate-new:
	@migrate create -ext sql -dir ${MIGRATION_PATH} -seq ${name}

test:
	@go test ./...

test-ci:
	@go test -v ./...

conduit-spec:
	APIURL=${API_URL} bash ./conduit/spec/run-api-tests.sh