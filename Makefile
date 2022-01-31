POSTGRES_URL="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
MIGRATION_PATH="./db/migrations"

migrate-up:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} up

migrate-down:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} down ${n}

migrate-force:
	@migrate -database ${POSTGRES_URL} -path ${MIGRATION_PATH} force ${v}

migrate-new:
	@migrate create -ext sql -dir ${MIGRATION_PATH} -seq ${name}