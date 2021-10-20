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

#### Description

Active miners are retrieved from Filecoin Lotus node.

First, the current active deals are retrieved.

Then active deals are parsed to get miners info, and store them in the database.

#### Diagram

<img src="./images/diagrams/get-miners.png" width="80%">

### Get locations

#### Description

Large airports are used to get relevant locations in the world.

They are imported from [https://datahub.io/core/airport-codes#data](https://datahub.io/core/airport-codes#data) and stored in the database.

#### Diagram

<img src="./images/diagrams/get-locations.png" width="30%">

### Get probes

#### Description

Get probes

#### Diagram

<img src="./images/diagrams/get-probes.png" width="80%">

### Create measures

#### Description

Create measures

#### Diagram

<img src="./images/diagrams/get-measurements.png" width="80%">

### Export measures

#### Description

Export reasults

#### Diagram

<img src="./images/diagrams/export-measurements.png" width="60%">

## JSON schema

[JSON Schema for exported data](./json/schema.json)

## Ripe Atlas costs
