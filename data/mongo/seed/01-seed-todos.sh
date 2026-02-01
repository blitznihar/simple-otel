#!/usr/bin/env bash
set -euo pipefail

echo "[mongo-init] Seeding todos from /docker-entrypoint-initdb.d/todos.json"

mongoimport \
  --db simple_otel \
  --collection todos \
  --jsonArray \
  --file /docker-entrypoint-initdb.d/todos.json

echo "[mongo-init] Done."