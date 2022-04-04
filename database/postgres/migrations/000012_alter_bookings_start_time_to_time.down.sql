ALTER TABLE bookings
    ALTER COLUMN start_time TYPE timestamp USING now()::date + start_time,
    ALTER COLUMN start_time SET NOT NULL,
    ALTER COLUMN end_time TYPE timestamp USING now()::date + end_time,
    ALTER COLUMN end_time SET NOT NULL;