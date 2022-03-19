ALTER TABLE places
    ALTER
        COLUMN open_hour TYPE TIME,
    ALTER
        COLUMN open_hour SET NOT NULL,
    ALTER
        COLUMN close_hour TYPE TIME,
    ALTER
        COLUMN close_hour SET NOT NULL;