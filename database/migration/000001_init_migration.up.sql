CREATE TABLE "users" (
  "uid" varchar PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "province_id" varchar,
  "income" varchar NOT NULL,
  "occupation_id" varchar,
  "photo" varchar NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
);

CREATE TABLE "occupations" (
  "occupation_id" varchar PRIMARY KEY,
  "title" varchar NOT NULL
);

CREATE TABLE "regencies" (
  "province_id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "alt_name" varchar NOT NULL,
  "latitude" varchar NOT NULL,
  "longitude" varchar NOT NULL
);

CREATE TABLE "expenses" (
  "uid" varchar UNIQUE NOT NULL,
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "type" varchar NOT NULL,
  "total" varchar NOT NULL,
  "Notes" varchar NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
);

ALTER TABLE "users" ADD FOREIGN KEY ("occupation_id") REFERENCES "occupations" ("occupation_id");

ALTER TABLE "users" ADD FOREIGN KEY ("province_id") REFERENCES "regencies" ("province_id");