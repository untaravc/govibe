-- +goose Up
CREATE TABLE IF NOT EXISTS `shipment_details` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `shipment_id` BIGINT UNSIGNED NOT NULL,
  `item_name` VARCHAR(191) NOT NULL,
  `item_price` DECIMAL(15,2) NOT NULL DEFAULT 0.00,
  `category_id` BIGINT UNSIGNED NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  KEY `shipment_details_shipment_id_index` (`shipment_id`),
  KEY `shipment_details_category_id_index` (`category_id`),
  KEY `shipment_details_deleted_at_index` (`deleted_at`),
  CONSTRAINT `shipment_details_shipment_id_foreign` FOREIGN KEY (`shipment_id`) REFERENCES `shipments` (`id`) ON DELETE CASCADE,
  CONSTRAINT `shipment_details_category_id_foreign` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS `shipment_details`;

