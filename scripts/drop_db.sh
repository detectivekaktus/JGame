#!/usr/bin/env bash
set -e

DBNAME=jgame_db
OWNER=jgame_owner
USER=jgame_backend

if [ $(whoami) != "postgres" ]; then
  echo "ERROR: You must run the script as postgres user."
  exit 1
fi

dropdb $DBNAME
dropuser $OWNER
dropuser $USER
