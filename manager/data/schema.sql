CREATE TABLE `locations` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text,`country` text,`iata_code` text,`latitude` real,`longitude` real,`type` text,PRIMARY KEY (`id`));
CREATE UNIQUE INDEX `idx_locations_iata_code` ON `locations`(`iata_code`);
CREATE INDEX `idx_locations_deleted_at` ON `locations`(`deleted_at`);
CREATE TABLE `locations_probes` (`location_id` integer,`probe_id` integer,PRIMARY KEY (`location_id`,`probe_id`),CONSTRAINT `fk_locations_probes_location` FOREIGN KEY (`location_id`) REFERENCES `locations`(`id`),CONSTRAINT `fk_locations_probes_probe` FOREIGN KEY (`probe_id`) REFERENCES `probes`(`id`));
CREATE INDEX locations_probes_location_id_index on locations_probes (location_id);
CREATE INDEX locations_probes_probe_id_index on locations_probes (probe_id);
CREATE TABLE `measurements` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`is_one_off` numeric,`measurement_id` integer,`start_time` integer,`stop_time` integer,`status` text,`status_stop_time` integer,PRIMARY KEY (`id`));
CREATE UNIQUE INDEX `idx_measurement_id` ON `measurements`(`measurement_id`);
CREATE INDEX `idx_measurements_deleted_at` ON `measurements`(`deleted_at`);
CREATE TABLE IF NOT EXISTS "miners"
(
	id integer
		primary key,
	created_at datetime,
	updated_at datetime,
	deleted_at datetime,
	address text,
	ip text,
	latitude real,
	longitude real,
	port integer
);
CREATE UNIQUE INDEX idx_address
	on miners (address);
CREATE INDEX idx_miners_deleted_at
	on miners (deleted_at);
CREATE TABLE IF NOT EXISTS "probes"
(
	id integer
		primary key,
	created_at datetime,
	updated_at datetime,
	deleted_at datetime,
	probe_id integer,
	country_code text,
	status text,
	latitude real,
	longitude real,
	is_anchor numeric,
	is_public numeric,
	coordinates_status text default "Unknown",
	ripe_latitude real,
	ripe_longitude real
, `address_v6` text, `address_v4` text);
CREATE INDEX idx_probes_coord_status
	on probes (coordinates_status);
CREATE INDEX idx_probes_deleted_at
	on probes (deleted_at);
CREATE UNIQUE INDEX idx_probes_probe_id
	on probes (probe_id);
CREATE INDEX idx_probes_status
	on probes (status);
CREATE TABLE `measurement_results` (`id` integer,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`probe_id` integer,`measurement_id` integer,`measurement_timestamp` integer,`ip` text,`measurement_date` text,`rtt` real,PRIMARY KEY (`id`),CONSTRAINT `fk_measurement_results_probe` FOREIGN KEY (`probe_id`) REFERENCES `probes`(`probe_id`),CONSTRAINT `fk_measurement_results_measurement` FOREIGN KEY (`measurement_id`) REFERENCES `measurements`(`measurement_id`));
CREATE INDEX `idx_measurement_date` ON `measurement_results`(`measurement_date`);
CREATE INDEX `idx_measurement_ip` ON `measurement_results`(`ip`);
CREATE INDEX `idx_measurement_results_deleted_at` ON `measurement_results`(`deleted_at`);
CREATE UNIQUE INDEX `idx_name` ON `measurement_results`(`probe_id`,`measurement_id`,`measurement_timestamp`,`ip`);
