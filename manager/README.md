# FC Latency Map - Manager

## Description

FC Latency Map - Manager allows to get latency measurements

## Requirements

### Packages

- sqlite3

```bash
sudo apt install sqlite3
```

### Configuration

Must have an ".env" file with the configurations to run.

Is available an example of configuration on [.env.example](./manager/.env.example)

```bash
cp .env.example .env
```

## Development

* Run golangci

```shell
# run golint-ci
golangci-lint run ./... --fix
```

* Install pre-commit tool

```shell
curl https://pre-commit.com/install-local.py | python -
```

* Install pre-commit hooks

```shell
# @ project root
pre-commit install
```

* Execute pre-commit manually

```shell
# @ project root
pre-commit run --all-files
```

## SQLite commands

1. Open database

```bash
sqlite3 data/database.db

SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite>
```

2. View tables

```bash
sqlite> .tables

>>> miners
```

3. Select miners

```bash
sqlite> select * from miners;

>>> 1|2021-09-07 17:04:53.41480159+02:00|2021-09-07 17:04:53.41480159+02:00||dummyAddress|dummyIp
```

4. Close database

```bash
sqlite> .quit
```
or
```bash
sqlite> ^C^C^C
```

## Command-line Interface
1. Run CLI
```bash
go run cmd/cli/main.go
```

2. Update miners list
without offset (it will apply by default the offset from the .env file)
```bash
>>> miners-update
```
or with an offset of 10 (latest block heights)
```bash
>>> miners-update 10
```

3. Parse miners from a given block height
```bash
>>> miners-parse 1107781
```

4. List miners
```bash
>>> miners-list
```


## Demo

```shell
## use case - from scratch
rm data/database.db

go run cmd/cli/main.go miners-parse 1109742
go run cmd/cli/main.go locations-add ORY
go run cmd/cli/main.go locations-add JFK
go run cmd/cli/main.go locations-add OPO

go run cmd/cli/main.go probes-update

go run cmd/cli/main.go measures-create
#   wait until have ripe results
go run cmd/cli/main.go measures-get
go run cmd/cli/main.go measures-export


## use case 2 - from seed data
#rm data/database.db
go run cmd/cli/main.go seed-data
go run cmd/cli/main.go probes-update

go run cmd/cli/main.go measures-get
go run cmd/cli/main.go measures-export
```
