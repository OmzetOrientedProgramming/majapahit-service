CREATE TABLE IF NOT EXISTS "disbursements" (
    "id" SERIAL PRIMARY KEY,
    "place_id" INT NOT NULL,
    "date" DATE NOT NULL, 
    "xendit_id" VARCHAR(128) NOT NULL,
    "amount" FLOAT NOT NULL,
    "status" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT now(),
    "updated_at" TIMESTAMP DEFAULT now(), 
    foreign key (place_id) references places(id)
);