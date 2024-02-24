CREATE TABLE `users` (
  `uid` varchar(255) PRIMARY KEY,
  `username` varchar(100) NOT NULL,
  `email` varchar(100) UNIQUE NOT NULL,
  `password` varchar(255) NOT NULL,
  `province_id` varchar(255),
  `occupation_id` varchar(255),
  `type_user` varchar(20) DEFAULT 'personal'
  `balance` varchar(15) DEFAULT '0',
  `savings` varchar(15) DEFAULT '0',
  `cash` varchar(15) DEFAULT '0',
  `debts` varchar(15) DEFAULT '0',
  `currency` varchar(10) DEFAULT 'IDR',
  `photo` varchar(255) DEFAULT '',
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
);

CREATE TABLE `occupations` (
  `occupation_id` varchar(255) PRIMARY KEY,
  `title` varchar(30) NOT NULL
);

CREATE TABLE `regencies` (
  `province_id` varchar(255) PRIMARY KEY,
  `name` varchar(20) NOT NULL,
  `alt_name` varchar(20) NOT NULL,
  `latitude` varchar(30) NOT NULL,
  `longitude` varchar(30) NOT NULL
);

CREATE TABLE `expenses` (
  `uid` varchar(255) UNIQUE NOT NULL,
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `type_x` varchar(10) NOT NULL,
  `total` varchar(15) NOT NULL,
  `Notes` varchar(100) DEFAULT '',
  `created_at` timestamp DEFAULT (now()),
  `updated_at` timestamp DEFAULT null
);

ALTER TABLE `users` ADD FOREIGN KEY (`occupation_id`) REFERENCES `occupations` (`occupation_id`);

ALTER TABLE `users` ADD FOREIGN KEY (`province_id`) REFERENCES `regencies` (`province_id`);
