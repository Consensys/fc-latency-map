                                                    |

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
