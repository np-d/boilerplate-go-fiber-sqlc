migrate-up:
	migrate -path ./sql/migration -database ${DATABASE_URL} -verbose up
migrate-down:
	migrate -path ./sql/migration -database ${DATABASE_URL} -verbose down