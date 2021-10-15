MIN_COV?=80

.PHONY: start coverage db-dump db-restore

start:
	go run cmd/cli/main.go

coverage:
	./scripts/coverage.sh $(MIN_COV)

db-dump:
	sqlite3 data/database.db .schema > data/schema.sql
	sqlite3 data/database.db .dump > data/dump.sql
	grep -vx -f data/schema.sql data/dump.sql > data/data.sql

db-restore:
	./scripts/db-restore.sh
