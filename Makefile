.PHONY: default manager map
.DEFAULT_GOAL := default

default: manager map
	mkdir -p exports
	docker-compose up

manager:
	docker build -f manager/Dockerfile -t fc-latency-manager .

map:
	docker build -f map/Dockerfile -t fc-latency-map ./map
