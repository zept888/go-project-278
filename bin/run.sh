#!/bin/sh
set -eu

echo "[run.sh] Starting service"

if [ -n "${DATABASE_URL:-}" ] && ls ./db/migrations/*.sql >/dev/null 2>&1; then
	echo "[run.sh] Running DB migrations"
	goose -dir ./db/migrations postgres "${DATABASE_URL}" up
elif [ -n "${DATABASE_URL:-}" ]; then
	echo "[run.sh] No migration files found, skipping migrations"
else
	echo "[run.sh] DATABASE_URL is not set, skipping migrations"
fi

echo "[run.sh] Starting Go API on :8081"
PORT=8081 /app/bin/app &
APP_PID=$!

trap 'kill "$APP_PID" 2>/dev/null || true' EXIT INT TERM

echo "[run.sh] Starting Caddy on :${PORT:-8080}"
exec caddy run --config /etc/caddy/Caddyfile --adapter caddyfile
