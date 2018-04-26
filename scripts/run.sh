#!/bin/sh
# Usage: run.sh SERVICE [CMD...]

# this is to ensure dependencies (ie postgres) get shutdown and rm'ed after run
# is completed
# https://github.com/docker/compose/issues/2791#issuecomment-198327517

docker-compose run -e TARGET=$TARGET --rm "$@"
exit_code=$?
docker-compose down -v
exit $exit_code
