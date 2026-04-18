-- +goose Up
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(191) NOT NULL,
  `email` VARCHAR(191) NOT NULL,
  `email_verified_at` DATETIME(3) NULL,
  `phone` VARCHAR(64) NULL,
  `phone_verified_at` DATETIME(3) NULL,
  `email_token` VARCHAR(255) NULL,
  `phone_token` VARCHAR(255) NULL,
  `image` VARCHAR(512) NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `role_id` BIGINT UNSIGNED NULL,
  `auth_type` ENUM('email','phone') NOT NULL DEFAULT 'email',
  `refresh_token` VARCHAR(255) NULL,
  `refresh_token_expired_at` DATETIME(3) NULL,
  `refresh_token_updated_at` DATETIME(3) NULL,
  `password` VARCHAR(255) NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_email_unique` (`email`),
  KEY `users_role_id_index` (`role_id`),
  KEY `users_deleted_at_index` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `users`
  ADD CONSTRAINT `users_role_id_foreign` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE SET NULL;

-- +goose Down
DROP TABLE IF EXISTS `users`;
