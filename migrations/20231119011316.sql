-- Create "channels" table
CREATE TABLE `channels` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` longtext NULL,
  `description` longtext NULL,
  `type` longtext NULL,
  `workspace_id` bigint unsigned NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workspaces_channels` (`workspace_id`),
  INDEX `idx_channels_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_workspaces_channels` FOREIGN KEY (`workspace_id`) REFERENCES `workspaces` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
