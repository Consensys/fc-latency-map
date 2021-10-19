# Filecoin Latency Map

## Config

Create an `.env` file:

```
cp .env.example .env
```

Config file values are:

| Key                   | Value type | Description              |
| --------------------- | ---------- | ------------------------ |
| SERVICE_NAME          | string     | The service name         |
| PATH_EXPORTS_MEASURES | string     | Json export directory    |
| LATENCY_LOW_LIMIT     | number     | Limit for low latency    |
| LATENCY_MEDIUM_LIMIT  | number     | Limit for medium latency |

## Run for development

Get packages:

```
npm i
```

Run the development server:

```bash
npm run dev
# or
yarn dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the map.

## Notification

To send exports notifications to the Map, send a Post request on [http://localhost:3000/api/notify](http://localhost:3000/api/notify)

It should respond:

```
{
"success": true
}
```

## Health Check

Open [http://localhost:3000/api/health-check](http://localhost:3000/api/health-check)

It should respond:

```
{
"success": true
}
```

## Metrics

Open [http://localhost:3000/api/metrics](http://localhost:3000/api/metrics)

It should respond:

```
{
"rssMemory": "1008.68 MB",
"heapMemoryTotal": "435.15 MB",
"heapMemoryUsed": "399.11 MB",
"externalMemory": "358.41 MB"
}
```

## Swagger

[./swagger.yml](./swagger.yml)
