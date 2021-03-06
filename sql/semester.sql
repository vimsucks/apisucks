CREATE TABLE IF NOT EXISTS `semesters` (
`id` INT UNSIGNED KEY PRIMARY KEY AUTO_INCREMENT,
`name` VARCHAR(15) NOT NULL,
`gpa` DOUBLE NOT NULL,
`year_id` INT UNSIGNED NOT NULL,
FOREIGN KEY (`year_id`) REFERENCES `years`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;