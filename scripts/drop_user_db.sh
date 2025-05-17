#!/usr/bin/env bash
set -e

DBNAME=jgamedb
USERNAME=jgameuser

if [ $(whoami) != "postgres" ]; then
  echo "ERROR: You must run the script as postgres user."
  exit 1
fi

dropdb $DBNAME
dropuser $USERNAME
