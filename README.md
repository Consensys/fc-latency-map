# FC Latency Map

[![CI pipeline](https://github.com/ConsenSys/fc-latency-map/actions/workflows/workflow.yml/badge.svg)](https://github.com/ConsenSys/fc-latency-map/actions/workflows/workflow.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ConsenSys/fc-latency-map)](https://goreportcard.com/report/github.com/ConsenSys/fc-latency-map)

[![MIT Licensed](https://img.shields.io/badge/License-MIT-brightgreen)](/LICENSE-MIT)
[![Apache Licensed](https://img.shields.io/badge/License-APACHE-brightgreen)](/LICENSE-APACHE)

## Description

FC Latency Map is a service for [Filecoin](https://filecoin.io/) blockchain to obtain the latencies of active miners.

It uses [Ripe Atlas](https://atlas.ripe.net/) to collect measurements of all active miners from a relevant location in the world.

## Quickstart

### Get the project

```shell
git clone https://github.com/ConsenSys/fc-latency-map.git
```

### Build the project

To build all the Docker images required to start the project, execute:

```shell
make
```

### Change default config

During build phase, `.env` config files were generated in `/manager` and `/map`. To start the services, 2 default values has to be changed on the manager config file.

Edit `/manager/.env` and change;

| Key               | Value type | Description             |
| ----------------- | ---------- | ----------------------- |
| FILECOIN_NODE_URL | string     | Lotus Filecoin node url |
| RIPE_API_KEY      | string     | Ripe Atlas API Key      |

### Start the project

Finally, to start the services, execute:

```shell
make run
```

The Fc Latency Map should be available at: [https://localhost:3000](https://localhost:3000)
