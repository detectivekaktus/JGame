#!/usr/bin/env bash
set -e

DBNAME=jgame_db
OWNER=jgame_owner
USER=jgame_backend

psql -U postgres -f database/00_init/01_create_roles.sql
psql -U postgres -f database/00_init/02_create_db.sql
psql -U $OWNER -d $DBNAME -f database/00_init/03_create_schemas.sql

for entry in database/*; do
  if [ -f "$entry" ] || [ $(basename "$entry") == "00_init" ]; then
    continue
  fi

  if [ -z "$( ls -A $entry )" ]; then
    echo "Skipping empty directory $entry"
    continue
  fi

  for file in $entry/*; do
    psql -U $OWNER -d $DBNAME -f $file
  done
done

