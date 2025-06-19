ALTER TABLE users."user"
ADD COLUMN matches_played INT DEFAULT 0,
ADD COLUMN matches_won INT DEFAULT 0;
