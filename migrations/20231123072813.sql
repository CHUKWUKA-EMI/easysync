-- Create "user_channels" table
CREATE TABLE `user_channels` (
  `user_id` char(36) NOT NULL,
  `channel_id` char(36) NOT NULL,
  PRIMARY KEY (`user_id`, `channel_id`),
  INDEX `fk_user_channels_channel` (`channel_id`),
  CONSTRAINT `fk_user_channels_channel` FOREIGN KEY (`channel_id`) REFERENCES `channels` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_user_channels_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
