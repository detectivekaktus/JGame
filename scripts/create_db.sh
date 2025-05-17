#!/usr/bin/env bash
# Creates a new database called jgame_db with owner jgame_owner.
# The owner role sets up the database for the user role called
# jgame_backend which, as the name suggest, is used on the backend.
# jgame_owner is never used after.
#
# The jgame_backend user follows the least privilege principle
# so they have only the necessary privileges to perform tasks on
# the tables inside the database schemas. The only deviation from
# the principle is the DELETE privilege given to the role, since
# the project isn't a real world application.
#
# The schemas created are:
# users:
#   user(user_id, email, name, password)
#   user_session(session_id, user_id, created_at, expires_at)
# packs:
#   pack(pack_id, user_id, body)
set -e

DBNAME=jgame_db
OWNER=jgame_owner
USER=jgame_backend

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <owner_password> <backend_password>"
  echo "Read the comment inside the script to see what it does."
  exit 1
fi

OWNER_PASSWORD="$1"
USER_PASSWORD="$2"

if [ $(whoami) != "postgres" ]; then
  echo "ERROR: You must run the script as postgres user."
  exit 1
fi

# TODO: Group the OWNER and the USER to one logical JGAME group.
psql -c "CREATE ROLE $OWNER WITH LOGIN PASSWORD '${OWNER_PASSWORD}' INHERIT;"
psql -c "CREATE ROLE $USER WITH LOGIN PASSWORD '${USER_PASSWORD}' INHERIT;"
psql -c "CREATE DATABASE $DBNAME OWNER $OWNER;"

psql -d $DBNAME -U $OWNER -c "CREATE SCHEMA users;"
psql -d $DBNAME -U $OWNER -c "CREATE SCHEMA packs;"

psql -d $DBNAME -U $OWNER -c "CREATE TABLE users.\"user\"(
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(32) NOT NULL,
  password TEXT NOT NULL
);"

psql -d $DBNAME -U $OWNER -c "CREATE TABLE users.user_session(
  session_id INT PRIMARY KEY,
  user_id INT,
  created_at TIMESTAMPTZ NOT NULL,
  expires_at TIMESTAMPTZ,
  FOREIGN KEY (user_id) REFERENCES users.\"user\"(user_id)
);"

psql -d $DBNAME -U $OWNER -c "CREATE TABLE packs.pack(
  pack_id SERIAL PRIMARY KEY,
  user_id INT,
  body JSONB NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users.\"user\"(user_id)
);"

psql -d $DBNAME -U $OWNER -c "REVOKE ALL ON SCHEMA public FROM PUBLIC;"
psql -d $DBNAME -U $OWNER -c "REVOKE ALL ON SCHEMA public FROM $USER;"

# TODO: Make future schemas inherit the defaults?
psql -d $DBNAME -U $OWNER -c "REVOKE CREATE ON SCHEMA users FROM $USER;"
psql -d $DBNAME -U $OWNER -c "REVOKE CREATE ON SCHEMA packs FROM $USER;"
psql -d $DBNAME -U $OWNER -c "GRANT USAGE ON SCHEMA users TO $USER;"
psql -d $DBNAME -U $OWNER -c "GRANT USAGE ON SCHEMA packs TO $USER;"

# NOTE: The DELETE privilege is given intentionally since this isn't a realworld
# hosted application where security and developer problems can have significant
# impact on the application.
psql -d $DBNAME -U $OWNER -c "GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA users TO $USER;"
psql -d $DBNAME -U $OWNER -c "GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA packs TO $USER;"
