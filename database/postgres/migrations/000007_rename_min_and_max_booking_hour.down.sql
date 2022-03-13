DO
$$
    BEGIN
        IF EXISTS(SELECT *
                  FROM information_schema.columns
                  WHERE table_name = 'places'
                    and column_name = 'min_interval_booking')
        THEN
            ALTER TABLE "public"."places"
                RENAME COLUMN "min_interval_booking" TO "min_hour_booking";
        END IF;
    END
$$;

DO
$$
    BEGIN
        IF
            EXISTS(SELECT *
                   FROM information_schema.columns
                   WHERE table_name = 'places'
                     and column_name = 'max_interval_booking')
        THEN
            ALTER TABLE "public"."places"
                RENAME COLUMN "max_interval_booking" TO "max_hour_booking";
        END IF;
    END
$$;