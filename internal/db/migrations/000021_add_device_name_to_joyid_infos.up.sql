ALTER TABLE joy_id_info_versions
    ADD old_device_name varchar(255) NOT NULL AFTER front_end,
    ADD device_name varchar(255) NOT NULL AFTER old_device_name,
    ADD old_device_type varchar(255) NOT NULL AFTER device_name,
    ADD device_type varchar(255) NOT NULL AFTER old_device_type;

ALTER TABLE joy_id_infos
    ADD device_name varchar(255) NOT NULL AFTER front_end,
    ADD device_type varchar(255) NOT NULL AFTER device_name;

ALTER TABLE sub_key_infos
    ADD device_name varchar(255) NOT NULL AFTER front_end,
    ADD device_type varchar(255) NOT NULL AFTER device_name;

CREATE TABLE IF NOT EXISTS sub_key_info_versions (
    id bigint NOT NULL AUTO_INCREMENT,
    old_block_number bigint unsigned NOT NULL,
    block_number bigint unsigned NOT NULL,
    lock_hash varchar(64) NOT NULL,
    pub_key char(128) NOT NULL,
    credential_id varchar(1500) NOT NULL,
    alg char(2) NOT NULL,
    old_front_end varchar(255) NOT NULL,
    front_end varchar(255) NOT NULL,
    old_device_name varchar(255) NOT NULL,
    device_name varchar(255) NOT NULL,
    old_device_type varchar(255) NOT NULL,
    device_type varchar(255) NOT NULL,
    action_type tinyint unsigned NOT NULL,
    tx_index int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_subkey_info_versions_on_block_number (block_number, action_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
