CREATE TABLE `accounts` (
  `user_id` varchar(255) NOT NULL,
  `account_id` varchar(255) UNIQUE NOT NULL,
  `account_name` varchar(50) NOT NULL DEFAULT '',
  `balance` int DEFAULT 0,
  `cash` int DEFAULT 0,
  `savings` int DEFAULT 0,
  `debts` int DEFAULT 0,
  `currency` varchar(10),
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT null ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE `accounts` ADD CONSTRAINT fk_accounts_users FOREIGN KEY (`user_id`) REFERENCES `users` (`uid`);