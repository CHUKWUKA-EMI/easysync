-- Modify "users" table
ALTER TABLE `users` ADD COLUMN `profile_state` enum('profile_active','profile_invited','profile_deactivated') NULL;
