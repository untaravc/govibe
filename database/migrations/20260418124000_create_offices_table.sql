-- +goose Up
CREATE TABLE IF NOT EXISTS `offices` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `type` VARCHAR(64) NOT NULL,
  `name` VARCHAR(191) NOT NULL,
  `code` VARCHAR(64) NOT NULL,
  `address` TEXT NULL,
  `phone` VARCHAR(64) NULL,
  `province_id` BIGINT UNSIGNED NULL,
  `city_id` BIGINT UNSIGNED NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `image_url` VARCHAR(512) NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `offices_code_unique` (`code`),
  KEY `offices_type_index` (`type`),
  KEY `offices_province_id_index` (`province_id`),
  KEY `offices_city_id_index` (`city_id`),
  KEY `offices_deleted_at_index` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS `offices`;

