CREATE TABLE IF NOT EXISTS "place_category" (
                                                "place_id" int not null,
                                                "category_id" int not null,
                                                foreign key (place_id) references places(id),
                                                foreign key (category_id) references categories(id)
);