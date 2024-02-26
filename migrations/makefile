HOST=localhost
DATABASE=veeektor_db
USERNAME=veeektor
PORT=5432

migrate:
	psql -h $(HOST) -d $(DATABASE) -U $(USERNAME) -p $(PORT) < ./create_tables.sql
