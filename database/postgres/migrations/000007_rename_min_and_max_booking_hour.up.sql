DO
$$
    BEGIN
        IF EXISTS(SELECT *
                  FROM information_schema.columns
                  WHERE table_name = 'places'
                    and column_name = 'min_hour_booking')
        THEN
            ALTER TABLE "public"."places"
                RENAME COLUMN "min_hour_booking" TO "min_interval_booking";
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
                     and column_name = 'max_hour_booking')
        THEN
            ALTER TABLE "public"."places"
                RENAME COLUMN "max_hour_booking" TO "max_interval_booking";
        END IF;
    END
$$;