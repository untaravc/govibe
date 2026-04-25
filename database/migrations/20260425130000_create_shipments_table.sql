-- +goose Up
CREATE TABLE IF NOT EXISTS `shipments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(64) NOT NULL,
  `customer_name` VARCHAR(191) NOT NULL,
  `office_origin_id` BIGINT UNSIGNED NOT NULL,
  `office_destination_id` BIGINT UNSIGNED NOT NULL,
  `customer_phone` VARCHAR(64) NOT NULL,
  `customer_email` VARCHAR(191) NULL,
  `price` DECIMAL(15,2) NOT NULL DEFAULT 0.00,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `wight` DECIMAL(12,3) NOT NULL DEFAULT 0.000,
  `length` DECIMAL(12,3) NOT NULL DEFAULT 0.000,
  `width` DECIMAL(12,3) NOT NULL DEFAULT 0.000,
  `height` DECIMAL(12,3) NOT NULL DEFAULT 0.000,
  `price_type` ENUM('dimension','weight') NOT NULL DEFAULT 'weight',
  `status` MEDIUMINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `shipments_code_unique` (`code`),
  KEY `shipments_office_origin_id_index` (`office_origin_id`),
  KEY `shipments_office_destination_id_index` (`office_destination_id`),
  KEY `shipments_user_id_index` (`user_id`),
  KEY `shipments_status_index` (`status`),
  KEY `shipments_deleted_at_index` (`deleted_at`),
  CONSTRAINT `shipments_office_origin_id_foreign` FOREIGN KEY (`office_origin_id`) REFERENCES `offices` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `shipments_office_destination_id_foreign` FOREIGN KEY (`office_destination_id`) REFERENCES `offices` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `shipments_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS `shipments`;
