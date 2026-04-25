-- +goose Up
CREATE TABLE IF NOT EXISTS `shipment_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `shipment_id` BIGINT UNSIGNED NOT NULL,
  `office_id` BIGINT UNSIGNED NOT NULL,
  `arrival_time` DATETIME(3) NULL,
  `departure_time` DATETIME(3) NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `note` TEXT NULL,
  `status` MEDIUMINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `shipment_logs_shipment_id_index` (`shipment_id`),
  KEY `shipment_logs_office_id_index` (`office_id`),
  KEY `shipment_logs_user_id_index` (`user_id`),
  KEY `shipment_logs_status_index` (`status`),
  KEY `shipment_logs_arrival_time_index` (`arrival_time`),
  KEY `shipment_logs_departure_time_index` (`departure_time`),
  KEY `shipment_logs_deleted_at_index` (`deleted_at`),
  CONSTRAINT `shipment_logs_shipment_id_foreign` FOREIGN KEY (`shipment_id`) REFERENCES `shipments` (`id`) ON DELETE CASCADE,
  CONSTRAINT `shipment_logs_office_id_foreign` FOREIGN KEY (`office_id`) REFERENCES `offices` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `shipment_logs_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS `shipment_logs`;

