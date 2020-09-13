SET NAMES utf8;
SET time_zone = '+03:00';
DROP TABLE IF EXISTS users;

CREATE TABLE `users` (
  `id` INT NOT NULL,
  `balance` DECIMAL(14, 2) NOT NULL,
   PRIMARY KEY (`id`)
)  ENGINE=InnoDB DEFAULT CHARSET=utf8;

