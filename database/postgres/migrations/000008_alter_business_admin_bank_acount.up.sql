DO
$$
    BEGIN
        IF EXISTS(SELECT *
                  FROM information_schema.columns
                  WHERE table_name = 'business_owners'
                    and column_name = 'bank_account')
        THEN
            ALTER TABLE "public"."business_owners"
                RENAME COLUMN "bank_account" TO "bank_account_number";
        END IF;
    END
$$;

ALTER TABLE business_owners
    ALTER COLUMN bank_account_number TYPE varchar(32);