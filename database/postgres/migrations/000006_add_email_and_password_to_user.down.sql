ALTER TABLE users
    DROP
        COLUMN IF EXISTS email,
    DROP
        COLUMN IF EXISTS password;