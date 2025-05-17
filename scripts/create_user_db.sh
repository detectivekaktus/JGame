#!/usr/bin/env bash
set -e

DBNAME=jgamedb
USERNAME=jgameuser

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <db_password>"
  echo "This script creates *$USERNAME* user with the provided password"
  echo "via the command line arguments. The user owns the *$DBNAME* database"
  echo "which is used to store the users, the packs, and the user sessions"
  echo "of the application. The user can perform CRUD operations on the database."
  exit 1
fi

PASSWORD="$1"

if [ $(whoami) != "postgres" ]; then
  echo "ERROR: You must run the script as postgres user."
  exit 1
fi

psql -c "CREATE ROLE $USERNAME WITH LOGIN PASSWORD '${PASSWORD}' INHERIT;"
psql -c "CREATE DATABASE $DBNAME OWNER $USERNAME;"

echo "Done creating *$DBNAME* database for *$USERNAME* user."
