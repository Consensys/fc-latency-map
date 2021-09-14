.PHONY: db-dump

db-dump:
	sqlite3 data/database.db .schema > data/schema.sql
	sqlite3 data/database.db .dump > data/dump.sql
	grep -vx -f data/schema.sql data/dump.sql > data/data.sql