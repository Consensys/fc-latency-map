# FC Latency Map Technical Design

## Architecture diagrams

### System

<img src="./images/fc-latency-map-architecture.png" width="800">

<strong>Ripe Atlas:</strong> open, distributed Internet measurement platform that measure Internet connectivity in real time.

<strong>Filecoin Node:</strong> node connected to Filecoin blockchain

<strong>Manager:</strong> Go service to create and export measures in JSON.

<strong>Export folder:</strong> JSON measurement results folder.

<strong>Map:</strong> React application to display measurements.

### Database

Database [dbdiagram model file](./filecoin_latency_map_dbdiagram)

<img src="./images/filecoin_latency_map_dbdiagram.png" width="800">

## Sequence diagrams

[Mermaid models file](./diagrams_mermaid.mmd)

### Get miners

Get miners

![get-miners](./images/diagrams/get-miners.png)

### Get locations

### Get probes

Get probes

![get-probes](./images/diagrams/get-probes.png)

### Create measures

Create measures

![get-measurements](./images/diagrams/get-measurements.png)

### Export measures

Export reasults

![export-measurements](./images/diagrams/export-measurements.png)

## JSON schema

[JSON Schema for exported data](./json/schema.json)

## Ripe Atlas costs
