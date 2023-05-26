CREATE TABLE `promotion`.`reusable_code` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(20),
    `is_active` TINYINT(4) DEFAULT '1',
    `max_use` INT NULL,
    `count` INT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=UTF8MB4 COLLATE=utf8mb4_unicode_ci;
