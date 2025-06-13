CREATE TABLE rooms.room(
  room_id INT PRIMARY KEY,
  user_id INT NOT NULL,
  name VARCHAR(32) NOT NULL,
  pack_id INT NOT NULL,
  current_users INT NOT NULL,
  max_users INT NOT NULL,
  password VARCHAR(32),
  FOREIGN KEY (pack_id) REFERENCES packs.pack(pack_id),
  FOREIGN KEY (user_id) REFERENCES users."user"(user_id)
);

CREATE TABLE rooms.player(
  user_id INT NOT NULL,
  room_id INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users."user"(user_id),
  FOREIGN KEY (room_id) REFERENCES rooms.room(room_id)
);
