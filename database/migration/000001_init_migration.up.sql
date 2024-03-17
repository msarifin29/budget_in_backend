CREATE TABLE `users` (
  `uid` varchar(255) PRIMARY KEY,
  `username` varchar(100) NOT NULL,
  `email` varchar(100) UNIQUE NOT NULL,
  `password` varchar(255) NOT NULL,
  `photo` varchar(255) DEFAULT '',
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);

-- CREATE TABLE `occupations` (
--   `occupation_id` varchar(255) PRIMARY KEY,
--   `title` varchar(30) NOT NULL,
--   `user_id` varchar(255) NOT NULL
-- );

-- CREATE TABLE `regencies` (
--   `province_id` varchar(255) PRIMARY KEY,
--   `name` varchar(20) NOT NULL,
--   `alt_name` varchar(20) NOT NULL,
--   `latitude` varchar(30) NOT NULL,
--   `longitude` varchar(30) NOT NULL,
--   `user_id` varchar(255) NOT NULL
-- );

CREATE TABLE `expenses` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `expense_type` varchar(10) NOT NULL,
  `total` INT NOT NULL DEFAULT 0,
  `Notes` text,
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT null ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE `incomes` (
  `uid` VARCHAR(255) NOT NULL,
  `id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `category_income` VARCHAR(15) not NULL DEFAULT 'monthly',
  `total` INT NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE `credits` (
  `uid` varchar(255) NOT NULL,
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `category_credit` varchar(50) DEFAULT 'other',
  `type_credit` varchar(15) NOT NULL DEFAULT 'monthly',
  `total` int NOT NULL,
  `loan_term` int NOT NULL DEFAULT 0,
  `status_credit` varchar(10) NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT null ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE `history_credit` (
  `credit_id` int NOT NULL,
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `th` int NOT NULL,
  `total` int NOT NULL,
  `status` varchar(10) NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT null ON UPDATE CURRENT_TIMESTAMP
);

-- ALTER TABLE regencies add CONSTRAINT fk_regencies_users FOREIGN KEY (user_id) REFERENCES users (uid)

-- ALTER TABLE occupations add CONSTRAINT fk_occupations_users FOREIGN KEY (user_id) REFERENCES users (uid)

ALTER TABLE `users` ADD COLUMN `type_user` varchar(10) DEFAULT 'personal';
ALTER TABLE `users` ADD COLUMN `balance` INT DEFAULT 0;
ALTER TABLE `users` ADD COLUMN `savings` INT DEFAULT 0;
ALTER TABLE `users` ADD COLUMN `cash` INT DEFAULT 0;
ALTER TABLE `users` ADD COLUMN `debts` INT DEFAULT 0;
ALTER TABLE `users` ADD COLUMN `currency` varchar(10) DEFAULT 'IDR';

ALTER TABLE `expenses` ADD COLUMN `uid` varchar(255) NOT NULL;
ALTER TABLE `expenses` ADD COLUMN `category` varchar(50) DEFAULT 'other';
ALTER TABLE `expenses` ADD COLUMN `status` varchar(10) DEFAULT 'success';
ALTER TABLE `expenses` ADD COLUMN `transaction_id` varchar(255) NOT NULL;

ALTER TABLE `incomes` ADD COLUMN `type_income` VARCHAR(15) not NULL DEFAULT 'debit';
ALTER TABLE `incomes` ADD COLUMN `transaction_id` varchar(255) NOT NULL;

ALTER TABLE `credits` ADD CONSTRAINT fk_credits_users FOREIGN KEY (`uid`) REFERENCES `users` (`uid`);

ALTER TABLE `history_credit` ADD CONSTRAINT fk_historycredit_credits FOREIGN KEY (`credit_id`) REFERENCES `credits` (`id`);

ALTER TABLE `credits` ADD COLUMN `installment` INT NOT NULL;
ALTER TABLE `credits` ADD COLUMN `payment_time` INT NOT NULL;

ALTER TABLE `history_credit` ADD COLUMN `type_payment` VARCHAR(10) not NULL DEFAULT 'cash';
ALTER TABLE `history_credit` ADD COLUMN `payment_time` int NOT NULL;