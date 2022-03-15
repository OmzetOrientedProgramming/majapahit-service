ALTER TABLE business_owners
    ADD COLUMN IF NOT EXISTS bank_account_name varchar(64) NOT NULL DEFAULT '';