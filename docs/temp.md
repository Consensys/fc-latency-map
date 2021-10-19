| Key                     | Value type | Description                                                                                               |
| ----------------------- | ---------- | --------------------------------------------------------------------------------------------------------- |
| RIPE_API_KEY            | string     | [Ripe Atlas key management](https://atlas.ripe.net/keys/)                                                 |
| RIPE_PING_INTERVAL      | number     | Interval between ping to get miners latency                                                               |
| RIPE_PING_RUNNING_TIME  | number     | Running period get latency                                                                                |
| RIPE_ONE_OFF            | boolean    | On ping only to get latency. When is 'true' the RIPE_PING_INTERVAL and RIPE_PING_RUNNING_TIME are ignored |
| RIPE_REQUESTED_PROBES   | number     | Max number of probes to call miners                                                                       |
| RIPE_PACKETS            | number     | Amount of packets to ping miners.                                                                         |
| NEAREST_AIRPORTS        | number     | Amount of airports miner to test latency.                                                                 |
| RIPE_PROBES_PER_AIRPORT | number     | Amount of probes near to miner to airport.                                                                |

Edit .env to add a valid Ripe Atlas API Key

After update the .env file execute:

```shell
go run cmd/cli/cli.go
```

### Diagrams

Get miners

![get-miners](./docs/diagrams/get-miners.png)

Get probes

![get-probes](./docs/diagrams/get-probes.png)

Create measurements

![get-measurements](./docs/diagrams/get-measurements.png)

Export reasults

![export-measurements](./docs/diagrams/export-measurements.png)

### Documentation

[./manager/README.md](./manager/README.md)

[JSON Schema for data](./docs/json/schema.json)

### License

The code is licensed under:

- [Apache v2.0 license](./LICENSE-APACHE).
- [MIT license](./LICENSE-MIT).
