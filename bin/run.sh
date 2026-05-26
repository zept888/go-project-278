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

echo "[run.sh] Starting Go app"
exec /app/bin/app
