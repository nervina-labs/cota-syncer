CREATE TABLE IF NOT EXISTS joy_id_infos (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    lock_hash varchar(64) NOT NULL,
    version varchar(40) NOT NULL,
    `name` varchar(255) NOT NULL,
    avatar varchar(500) NOT NULL,
    description varchar(1000) NOT NULL,
    extension varchar(1000) NOT NULL,
    pub_key char(128) NOT NULL,
    credential_id varchar(100) NOT NULL,
    alg char(2) NOT NULL,
    cota_cell_id char(16) NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_joy_id_infos_on_block_number (block_number),
    CONSTRAINT uc_joy_id_infos_on_lock_hash UNIQUE (lock_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS joy_id_info_versions (
    id bigint NOT NULL AUTO_INCREMENT,
    old_block_number bigint unsigned NOT NULL,
    block_number bigint unsigned NOT NULL,
    lock_hash varchar(64) NOT NULL,
    old_version varchar(40) NOT NULL,
    version varchar(40) NOT NULL,
    old_name varchar(255) NOT NULL,
    `name` varchar(255) NOT NULL,
    old_avatar varchar(500) NOT NULL,
    avatar varchar(500) NOT NULL,
    old_description varchar(1000) NOT NULL,
    description varchar(1000) NOT NULL,
    old_extension varchar(1000) NOT NULL,
    extension varchar(1000) NOT NULL,
    pub_key char(128) NOT NULL,
    credential_id varchar(100) NOT NULL,
    alg char(2) NOT NULL,
    cota_cell_id char(16) NOT NULL,
    action_type tinyint unsigned NOT NULL,
    tx_index int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_joy_id_versions_on_block_number (block_number, action_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE IF NOT EXISTS sub_key_infos (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    lock_hash varchar(64) NOT NULL,
    pub_key char(128) NOT NULL,
    credential_id varchar(100) NOT NULL,
    alg char(2) NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_sub_key_infos_on_block_number (block_number),
    CONSTRAINT uc_sub_key_infos_on_pub_key UNIQUE (pub_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;