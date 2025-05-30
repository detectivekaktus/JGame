GRANT USAGE ON SCHEMA users TO jgame_backend;
GRANT USAGE ON SCHEMA packs TO jgame_backend;

ALTER DEFAULT PRIVILEGES IN SCHEMA users
  GRANT SELECT, INSERT, UPDATE, DELETE
  ON TABLES TO jgame_backend;
ALTER DEFAULT PRIVILEGES IN SCHEMA users
  GRANT USAGE, SELECT ON SEQUENCES TO jgame_backend;

ALTER DEFAULT PRIVILEGES IN SCHEMA packs 
  GRANT SELECT, INSERT, UPDATE, DELETE
  ON TABLES TO jgame_backend;
ALTER DEFAULT PRIVILEGES IN SCHEMA packs
  GRANT USAGE, SELECT ON SEQUENCES TO jgame_backend;
