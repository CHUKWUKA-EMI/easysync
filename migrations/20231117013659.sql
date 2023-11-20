-- Modify "users" table
ALTER TABLE `users` DROP COLUMN `on_boarding_step`, ADD COLUMN `onboarding_step` longtext NULL;
