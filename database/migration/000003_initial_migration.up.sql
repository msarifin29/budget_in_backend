CREATE TABLE "users" (
  "uid" varchar(255) PRIMARY KEY,
  "username" varchar(100) NOT NULL,
  "email" varchar(100) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL,
  "photo" varchar(255) DEFAULT '',
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT null,
  "type_user" varchar(10) DEFAULT 'personal',
  "status" varchar(15) DEFAULT 'active'
);

CREATE TABLE "expenses" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "expense_type" varchar(10) NOT NULL,
  "total" INT NOT NULL DEFAULT 0,
  "Notes" text,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT null,
  "uid" varchar(255) NOT NULL,
  "status" varchar(10) DEFAULT 'success',
  "transaction_id" varchar(255) NOT NULL
);

CREATE TABLE "incomes" (
  "uid" VARCHAR(255) NOT NULL,
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY NOT NULL,
  "total" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT null,
  "type_income" VARCHAR(15) NOT NULL DEFAULT 'debit',
  "transaction_id" varchar(255) NOT NULL
);

CREATE TABLE "t_category_expenses" (
  "category_id" int UNIQUE PRIMARY KEY NOT NULL,
  "id" int NOT NULL,
  "title" varchar(50) NOT NULL
);

CREATE TABLE "t_category_incomes" (
  "category_id" int UNIQUE PRIMARY KEY NOT NULL,
  "id" int NOT NULL,
  "title" varchar(50) NOT NULL
);

CREATE TABLE "accounts" (
  "user_id" varchar(255) NOT NULL,
  "account_id" varchar(255) UNIQUE NOT NULL,
  "account_name" varchar(50) NOT NULL DEFAULT '',
  "balance" int DEFAULT 0,
  "cash" int DEFAULT 0,
  "savings" int DEFAULT 0,
  "debts" int DEFAULT 0,
  "currency" varchar(10),
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT null,
  "max_budget" int DEFAULT 0
);

ALTER TABLE "t_category_expenses" ADD CONSTRAINT "fk_t_category_expenses_expenses" FOREIGN KEY ("category_id") REFERENCES "expenses" ("id");

ALTER TABLE "t_category_incomes" ADD CONSTRAINT "fk_t_category_incomes_incomes" FOREIGN KEY ("category_id") REFERENCES "incomes" ("id");

ALTER TABLE "accounts" ADD CONSTRAINT "fk_accounts_users" FOREIGN KEY ("user_id") REFERENCES "users" ("uid");
