CREATE TABLE rooms.room(
  room_id TEXT PRIMARY KEY,
  name VARCHAR(32) NOT NULL,
  pack_id INT NOT NULL,
  current_users INT NOT NULL,
  max_users INT NOT NULL,
  password VARCHAR(32),
  FOREIGN KEY (pack_id) REFERENCES packs.pack(pack_id)
);

CREATE TABLE rooms.players(
  user_id INT NOT NULL,
  room_id TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users."user"(user_id),
  FOREIGN KEY (room_id) REFERENCES rooms.room(room_id)
);
