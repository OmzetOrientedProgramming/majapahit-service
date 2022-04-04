CREATE TABLE IF NOT EXISTS "time_slots"
(
    "id"         serial primary key,
    "start_time" time not null,
    "end_time"   time not null,
    "day"        int  not null,
    "place_id"   int references places (id),
    "created_at" timestamp default now(),
    "updated_at" timestamp default now()
);