-- Create "workspaces" table
CREATE TABLE `workspaces` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` longtext NULL,
  `logo_url` longtext NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_workspaces_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "channels" table
CREATE TABLE `channels` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` longtext NULL,
  `description` longtext NULL,
  `type` longtext NULL,
  `owner_email` longtext NULL,
  `workspace_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_workspaces_channels` (`workspace_id`),
  INDEX `idx_channels_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_workspaces_channels` FOREIGN KEY (`workspace_id`) REFERENCES `workspaces` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "roles" table
CREATE TABLE `roles` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` longtext NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_roles_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "users" table
CREATE TABLE `users` (
  `id` char(36) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `email` varchar(191) NOT NULL,
  `is_email_confirmed` bool NULL DEFAULT 0,
  `password` longtext NULL,
  `first_name` longtext NULL,
  `last_name` longtext NULL,
  `real_name` longtext NULL,
  `display_name` longtext NULL,
  `occupation` longtext NULL,
  `profile_image_url` longtext NULL,
  `onboarding_step` longtext NULL,
  `profile_state` enum('profile_active','profile_invited','profile_deactivated') NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_users_deleted_at` (`deleted_at`),
  INDEX `idx_users_email` (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "user_roles" table
CREATE TABLE `user_roles` (
  `user_id` char(36) NOT NULL,
  `role_id` char(36) NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`),
  INDEX `fk_user_roles_role` (`role_id`),
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "user_workspaces" table
CREATE TABLE `user_workspaces` (
  `user_id` char(36) NOT NULL,
  `workspace_id` char(36) NOT NULL,
  PRIMARY KEY (`user_id`, `workspace_id`),
  INDEX `fk_user_workspaces_workspace` (`workspace_id`),
  CONSTRAINT `fk_user_workspaces_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_user_workspaces_workspace` FOREIGN KEY (`workspace_id`) REFERENCES `workspaces` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
