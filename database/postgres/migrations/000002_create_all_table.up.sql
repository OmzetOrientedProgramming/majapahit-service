CREATE TABLE IF NOT EXISTS "users" (
    "id" serial primary key,
    "phone_number" varchar(15) unique not null,
    "name" varchar(50) not null,
    "status" int not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now()
);

CREATE TABLE IF NOT EXISTS "business_owners" (
    "id" serial primary key,
    "balance" float not null,
    "bank_account" int not null,
    "user_id" int not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now(),
    foreign key (user_id) references users(id)
);

CREATE TABLE IF NOT EXISTS "customers" (
    "id" serial primary key,
    "date_of_birth" date not null ,
    "gender" int not null,
    "user_id" int not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now(),
    foreign key (user_id) references users(id)
);

CREATE TABLE IF NOT EXISTS "places" (
     "id" serial primary key,
     "name" varchar(50) not null,
     "address" varchar(100) not null,
     "capacity" int not null,
     "description" text not null,
     "user_id" int not null,
     "interval" int not null,
     "open_hour" timestamp not null,
     "close_hour" timestamp not null,
     "image" varchar(50) not null,
     "min_hour_booking" int not null,
     "max_hour_booking" int not null,
     "min_slot_booking" int not null,
     "max_slot_booking" int not null,
     "lat" float not null,
     "long" float not null,
     "created_at" timestamp default now(),
     "updated_at" timestamp default now(),
     foreign key (user_id) references users(id)
);

CREATE TABLE IF NOT EXISTS "bookings" (
    "id" serial primary key,
    "user_id" int not null,
    "place_id" int not null,
    "date" date not null,
    "start_time" timestamp not null,
    "end_time" timestamp not null,
    "capacity" int not null,
    "status" int not null,
    "total_price" float not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now(),
    foreign key (user_id) references users(id),
    foreign key (place_id) references places(id)
);

CREATE TABLE IF NOT EXISTS "reviews" (
    "id" serial primary key,
    "user_id" int not null,
    "place_id" int not null,
    "booking_id" int not null,
    "content" text not null,
    "rating" int not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now(),
    foreign key (user_id) references users(id),
    foreign key (place_id) references places(id)
);

CREATE TABLE IF NOT EXISTS "items" (
    "id" serial primary key,
    "name" varchar(15) not null,
    "image" varchar(50) not null,
    "price" float not null,
    "description" text not null,
    "place_id" int not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now(),
    foreign key (place_id) references places(id)
);

CREATE TABLE IF NOT EXISTS "booking_items" (
    "id" serial primary key,
    "item_id" int not null,
    "booking_id" int not null,
    "qty" int not null,
    "total_price" float not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now(),
    foreign key (item_id) references items(id),
    foreign key (booking_id) references bookings(id)
);


CREATE TABLE IF NOT EXISTS "tags" (
    "id" int primary key,
    "name" varchar(15) not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now()
);


CREATE TABLE IF NOT EXISTS  "place_tags" (
    "id" serial primary key,
    "tag_id" int not null,
    "place_id" int not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now(),
    foreign key (tag_id) references tags(id),
    foreign key (place_id) references places(id)
);

CREATE TABLE IF NOT EXISTS "otp_codes" (
    "id" serial primary  key,
    "otp" int not null,
    "expired_date" timestamp not null ,
    "phone_number" varchar(15) not null,
    "created_at" timestamp default now(),
    "updated_at" timestamp default now()
);



















