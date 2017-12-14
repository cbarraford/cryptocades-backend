#!/bin/sh
# wait-for-postgres.sh
# https://docs.docker.com/compose/startup-order/

set -e

echo "Waiting for Postgres..."

cmd="$@"

until /usr/bin/pg_isready; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd
