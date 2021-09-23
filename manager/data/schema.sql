create table locations
(
    id integer
        primary key,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime,
    country text,
    iata_code text,
    latitude real,
    longitude real
);

create index idx_locations_deleted_at
	on locations (deleted_at);

create index idx_locations_iata_code
	on locations (iata_code);

create table measurement_results
(
    id integer
        primary key,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime,
    probe_id integer,
    measurement_id integer,
    measurement_timestamp integer,
    measurement_date text,
    time_average real,
    time_max real,
    time_min real,
    ip text
);

create index idx_measurement_results_deleted_at
	on measurement_results (deleted_at);

create index idx_measurement_results_ip
	on measurement_results (ip);

create index idx_measurement_results_measurement_id
	on measurement_results (measurement_id);

create index idx_measurement_results_measurement_timestamp
	on measurement_results (measurement_timestamp);

create table measurements
(
    id integer
        primary key,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime,
    is_one_off numeric,
    measurement_id integer,
    start_time integer,
    stop_time integer
);

create index idx_measurements_deleted_at
	on measurements (deleted_at);

create table migrations
(
    id VARCHAR(255)
        primary key
);

create table miners
(
    id integer
        primary key,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime,
    address text,
    ip text,
    latitude real,
    longitude real
);

create unique index idx_miners_address
	on miners (address);

create index idx_miners_deleted_at
	on miners (deleted_at);

create table probes
(
    id integer
        primary key,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime,
    probe_id integer,
    country_code text,
    latitude real,
    longitude real
);

create index idx_probes_deleted_at
	on probes (deleted_at);

create unique index idx_probes_probe_id
	on probes (probe_id);

create table sqlite_master
(
    type text,
    name text,
    tbl_name text,
    rootpage int,
    sql text
);