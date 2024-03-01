CREATE TABLE `users` (
  `uid` varchar(255) PRIMARY KEY,
  `username` varchar(100) NOT NULL,
  `email` varchar(100) UNIQUE NOT NULL,
  `password` varchar(255) NOT NULL,
  `type_user` varchar(20) DEFAULT 'personal'
  `balance` INT DEFAULT 0,
  `savings` INT DEFAULT 0,
  `cash` INT DEFAULT 0,
  `debts` INT DEFAULT 0,
  `currency` varchar(10) DEFAULT 'IDR',
  `photo` varchar(255) DEFAULT '',
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE `occupations` (
  `occupation_id` varchar(255) PRIMARY KEY,
  `title` varchar(30) NOT NULL,
  `user_id` varchar(255) NOT NULL
);

CREATE TABLE `regencies` (
  `province_id` varchar(255) PRIMARY KEY,
  `name` varchar(20) NOT NULL,
  `alt_name` varchar(20) NOT NULL,
  `latitude` varchar(30) NOT NULL,
  `longitude` varchar(30) NOT NULL,
  `user_id` varchar(255) NOT NULL
);

CREATE TABLE `expenses` (
  `uid` varchar(255) NOT NULL,
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `expense_type` varchar(10) NOT NULL,
  `total` INT NOT NULL DEFAULT 0,
  `category` varchar(50) DEFAULT 'other',
  `status` varchar(10) DEFAULT 'success',
  `Notes` text,
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT null ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE `incomes` (
  `uid` VARCHAR(255) NOT NULL UNIQUE,
  `id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `category_income` VARCHAR(15) not NULL DEFAULT 'monthly',
  `total` INT NOT NULL DEFAULT 0,
  `type_income` VARCHAR(15) not NULL DEFAULT 'debit',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE regencies add CONSTRAINT fk_regencies_users FOREIGN KEY (user_id) REFERENCES users (uid)

ALTER TABLE occupations add CONSTRAINT fk_occupations_users FOREIGN KEY (user_id) REFERENCES users (uid)
