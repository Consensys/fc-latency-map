.PHONY: default run manager map
.DEFAULT_GOAL := default

default: 
	./scripts/db-restore.sh
	./scripts/generate-config.sh
	./scripts/build-manager.sh
	docker build -f map/Dockerfile -t fc-latency-map ./map

run:
	./run.sh

manager:
	docker build -f manager/Dockerfile -t fc-latency-manager .

map:
	docker build -f map/Dockerfile -t fc-latency-map ./map
