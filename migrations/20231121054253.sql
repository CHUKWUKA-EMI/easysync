-- Create "invites" table
CREATE TABLE `invites` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `token` longtext NULL,
  `status` enum('invite_pending','invite_accepted') NULL,
  `workspace_id` char(36) NULL,
  `user_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_invites_deleted_at` (`deleted_at`),
  INDEX `idx_member` (`user_id`, `workspace_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tokens" table
CREATE TABLE `tokens` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `refresh_token` longtext NULL,
  `user_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_tokens_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
