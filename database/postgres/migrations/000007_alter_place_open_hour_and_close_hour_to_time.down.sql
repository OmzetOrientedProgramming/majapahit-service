ALTER TABLE places
    ALTER
        COLUMN open_hour TYPE timestamp USING now()::date + open_hour,
    ALTER
        COLUMN open_hour SET NOT NULL,
    ALTER
        COLUMN close_hour TYPE timestamp USING now()::date + close_hour,
    ALTER
        COLUMN close_hour SET NOT NULL;