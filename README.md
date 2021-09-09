# FC Latency Map

## Description

FC Latency Map are services for [Filecoin](https://filecoin.io/) blockchain to get latencies of active miners.

It uses [Ripe Atlas](https://atlas.ripe.net/) to collect measurements of every active miners from relevant location on the world.

Get the project:

```
git clone https://github.com/ConsenSys/fc-latency-map.git
```

## Manager

### Quickstart

Start the Manager cli:

```
cd manager
cp .env.example .env
```

Edit .env to add a valid Ripe Atlas API Key, then execute:

```
go run cmd/cli/cli.go
```

### Documentation

[./manager/README.md](./manager/README.md)
