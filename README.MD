# simple-otel

[![CI](https://github.com/blitznihar/simple-otel/actions/workflows/ci.yml/badge.svg)](https://github.com/blitznihar/simple-otel/actions)
[![License](https://img.shields.io/github/license/blitznihar/simple-otel.svg)](LICENSE)

This is a perfect “trace a request across 5 services” playground for OpenTelemetry + Jaeger.

This repository is an opinionated blueprint that keeps tracing, API contracts, and instrumentation consistent across multiple languages and services: Python → Spring (Java) → Node → Go → MongoDB, with an Angular UI as the caller.

**Target flow (the “trace story”)**

Angular UI
→ calls Python API (edge/gateway)
→ Python calls Spring API
→ Spring calls Node API
→ Node calls Go API (persistence service)
→ Go reads/writes MongoDB

What this ensures:

- Every hop propagates the same `traceparent` header
- Each service exports spans via OTLP to an `otel-collector`
- Collector sends traces to Jaeger (for visualization)

**Standard API contract (same in every service)**

Use one contract everywhere to avoid model drift. Example `Todo` JSON:

```json
{
	"id": 1,
	"title": "Do dishes",
	"note": "That will be done by Gopher.",
	"due_date": "2026-02-15T00:00:00Z",
	"created_at": "2026-01-20T09:00:00Z",
	"is_complete": false
}
```

Endpoints (same in every service):

- `GET /health`
- `GET /todos` → returns list
- `GET /todos/{id}` → returns one
- `POST /todos` → create (returns created)
- `PUT /todos/{id}` → update (returns updated)

**Call downstream rule (important for tracing)**

For learning, each service should:

- Start a span for the inbound request (auto-instrumentation)
- Start a span for the outbound downstream call (HTTP client span)
- Add 2–4 attributes so traces are readable: `service.name`, `todo.id`, `downstream.service`, `operation` (create/update/get)

This makes traces readable and helps learners follow the request as it hops between services.

**Subprojects / Examples**

- [Python API](todo-api-py/README.md)
- [Node.js API](todo-api-nodejs/README.md)
- [Go API](todo-api-go/README.MD)
- [Spring API](todo-api-spring/README.MD)
- [Angular UI (caller)](todo-ng-ui/README.MD)

See the subproject READMEs for implementation details, instrumentation tips, and run instructions for each service.

License: See the [LICENSE](LICENSE) file.

