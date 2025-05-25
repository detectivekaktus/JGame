DO $$
BEGIN
  CREATE TABLE users."user"(
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(32) NOT NULL,
    password TEXT NOT NULL
  );

  CREATE TABLE users.user_session(
    session_id TEXT PRIMARY KEY,
    user_id INT,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users."user"(user_id)
  );

  CREATE TABLE packs.pack(
    pack_id SERIAL PRIMARY KEY,
    user_id INT,
    body JSONB NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users."user"(user_id)
  );
END
$$;
