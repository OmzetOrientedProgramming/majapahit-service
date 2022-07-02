CREATE TABLE IF NOT EXISTS "categories" (
                                       "id" serial primary key,
                                       "created_at" timestamp default now(),
                                       "updated_at" timestamp default now(),
                                       "content" varchar(255) not null
);