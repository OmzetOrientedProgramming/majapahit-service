ALTER TABLE bookings
    ALTER COLUMN start_time TYPE TIME,
    ALTER COLUMN start_time SET NOT NULL,
    ALTER COLUMN end_time TYPE TIME,
    ALTER COLUMN end_time SET NOT NULL;