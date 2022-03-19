ALTER TABLE users
    ADD COLUMN IF NOT EXISTS email    varchar(64) UNIQUE,
    ADD COLUMN IF NOT EXISTS password varchar(256);