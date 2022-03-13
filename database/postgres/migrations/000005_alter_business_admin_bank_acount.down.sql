DO
$$
    BEGIN
        IF EXISTS(SELECT *
                  FROM information_schema.columns
                  WHERE table_name = 'business_owners'
                    and column_name = 'bank_account_number')
        THEN
            ALTER TABLE "public"."business_owners"
                RENAME COLUMN "bank_account_number" TO "bank_account";
        END IF;
    END
$$;

ALTER TABLE business_owners
    ALTER COLUMN bank_account TYPE integer USING bank_account::integer;