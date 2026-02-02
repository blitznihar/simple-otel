# simple-otel

[![CI](https://github.com/blitznihar/simple-otel/actions/workflows/ci.yml/badge.svg)](https://github.com/blitznihar/simple-otel/actions)
[![License](https://img.shields.io/github/license/blitznihar/simple-otel.svg)](LICENSE)

This repository demonstrates instrumentation and end-to-end observability using OpenTelemetry (Jaeger) across simple TODO APIs implemented in Node.js and Go, with MongoDB as the backing store. It includes container and Kubernetes manifests so you can run locally (Docker) or on a cluster (Kubernetes).

## Contents

- `todo-api-nodejs` — Node.js example service (Express). Includes OTEL setup and sample routes.
- `todo-api-go` — Go example service. Includes OTEL setup and sample routes.
- `todo-api-py` — Python example (minimal).
- `data/mongo` — MongoDB seed data and init scripts.
- `infra/otel` — Kubernetes manifests for Jaeger and an OpenTelemetry Collector configuration.

## What is done

- Instrumentation: both the Node.js and Go example services include OpenTelemetry bootstrap files and are configured to export traces to Jaeger (or OTEL Collector) when appropriate environment variables are set.
- Persistence: a MongoDB instance is used for todos; seed data is provided in `data/mongo/seed/todos.json` and the `01-seed-todos.sh` script.
- Containerization: each service has a `Dockerfile` so you can run services in containers.
- Kubernetes: example manifests are provided under each service `k8s/` directory and `infra/otel` for tracing stack (Jaeger + OTEL collector).

## What is proven

- End-to-end traces: when Jaeger is running and services are configured to export telemetry, you can observe traces for incoming HTTP requests that include spans from the application and any outbound calls.
- Data flow: CRUD requests to the TODO APIs persist to MongoDB and can be verified via the provided seed and API endpoints.

## How to run locally (quick start)

1. Start MongoDB (Docker):

```bash
docker run -d --name simple-otel-mongo -p 27017:27017 mongo:6.0
```

2. Seed the database (optional):

```bash
bash data/mongo/seed/01-seed-todos.sh
```

3. Start Jaeger (all-in-one) so you can view traces at `http://localhost:16686`:

```bash
docker run --rm -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp \
  -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 14250:14250 \
  jaegertracing/all-in-one:1.47
```

4. Run the Node.js service (local):

```bash
cd todo-api-nodejs
npm install
# point OTEL exporter to local Jaeger - environment depends on the service implementation
export JAEGER_AGENT_HOST=localhost
npm start
```

5. Run the Go service (local):

```bash
cd todo-api-go
go run ./cmd
```

Notes: if you prefer the OTEL Collector, run the collector and set `OTEL_EXPORTER_OTLP_ENDPOINT` accordingly.

## How to test

- Verify API is up:

```bash
curl -sS http://localhost:3000/health || curl -sS http://localhost:8080/health
```

- List todos (example endpoints):

```bash
curl -sS http://localhost:3000/todos
curl -sS http://localhost:8080/todos
```

- Create a todo:

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"title":"test","completed":false}' \
  http://localhost:3000/todos
```

- Seeded data check:

```bash
mongo --eval 'db.todos.find().limit(5).pretty()' localhost:27017
```

## Verify tracing in Jaeger

1. Open the Jaeger UI at: `http://localhost:16686`.
2. Select the service name (e.g., `todo-api-nodejs` or `todo-api-go`) and search for traces.
3. Make several API requests (list/create) and refresh the Jaeger UI — spans and traces should appear showing the HTTP server span and any DB or client spans (if instrumented).

## Kubernetes

- The repository includes Kubernetes manifests for running the services and tracing stack under `infra/otel` and each service's `k8s/` directory. Apply them in a cluster with `kubectl apply -f <manifest>`.

## Troubleshooting

- No traces visible: verify your service environment variables point to the Jaeger agent/collector (`JAEGER_AGENT_HOST`, `OTEL_EXPORTER_OTLP_ENDPOINT`, etc.) and that ports (6831 UDP for Jaeger agent, or OTLP port) are reachable.
- Mongo connection issues: ensure Mongo is running on the expected host:port and that `todo-api-*.env` files (or k8s configmaps) are configured properly.

## Next steps

- I can run the services and verify traces end-to-end, or update per-service README sections with exact env vars used by each implementation. Which would you like me to do next?

License: See the [LICENSE](LICENSE) file.

