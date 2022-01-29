migrate-up:
	@migrate -database ${db} -path ./db/migrations up

migrate-down:
	@migrate -database ${db} -path ./db/migrations down ${n}

migrate-new:
	@migrate create -ext sql -dir ./db/migrations -seq ${name}