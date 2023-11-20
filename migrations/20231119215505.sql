-- Modify "users" table
ALTER TABLE `users` MODIFY COLUMN `profile_state` enum('profile_active','profile_invited','profile_deactivated') NULL;
