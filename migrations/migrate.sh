export POSTGRESQL_URL="postgresql://veeektor:766180@localhost:5432/veeektor_db?sslmode=disable"
migrate -database ${POSTGRESQL_URL} -path . up  