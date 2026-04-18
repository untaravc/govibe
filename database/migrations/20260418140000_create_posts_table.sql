-- +goose Up
CREATE TABLE IF NOT EXISTS `posts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(191) NOT NULL,
  `subtitle` VARCHAR(255) NULL,
  `slug` VARCHAR(191) NOT NULL,
  `content` LONGTEXT NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `image_url` VARCHAR(512) NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `release_at` DATETIME(3) NULL,
  `category_id` BIGINT UNSIGNED NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `posts_slug_unique` (`slug`),
  KEY `posts_user_id_index` (`user_id`),
  KEY `posts_category_id_index` (`category_id`),
  KEY `posts_release_at_index` (`release_at`),
  KEY `posts_deleted_at_index` (`deleted_at`),
  CONSTRAINT `posts_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS `posts`;

