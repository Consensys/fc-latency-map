CREATE TABLE `miners` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`address` text UNIQUE,`ip` text,PRIMARY KEY (`id`));
CREATE INDEX `idx_miners_deleted_at` ON `miners`(`deleted_at`);
CREATE TABLE `locations` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`country` text,`latitude` text,`longitude` text,PRIMARY KEY (`id`));
CREATE INDEX `idx_locations_deleted_at` ON `locations`(`deleted_at`);
CREATE TABLE `measurements` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`miner` text,`probe_id` integer,`measure_date` integer,`time_average` real,PRIMARY KEY (`id`));
CREATE INDEX `idx_measurements_deleted_at` ON `measurements`(`deleted_at`);
