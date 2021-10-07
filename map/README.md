## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the map.

## Health Check

Open [http://localhost:3000/api/health-check](http://localhost:3000/api/health-check)

Example:

```
{
"success": true
}
```

## Metrics

Open [http://localhost:3000/api/metrics](http://localhost:3000/api/metrics)

Example:

```
{
"rssMemory": "1008.68 MB",
"heapMemoryTotal": "435.15 MB",
"heapMemoryUsed": "399.11 MB",
"externalMemory": "358.41 MB"
}
```
