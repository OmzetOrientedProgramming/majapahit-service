ALTER TABLE users
    ADD COLUMN IF NOT EXISTS firebase_local_id varchar(65) UNIQUE;