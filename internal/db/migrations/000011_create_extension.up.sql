CREATE TABLE IF NOT EXISTS extension_kv_pairs (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    lock_hash char(64) NOT NULL,
    lock_hash_crc int unsigned NOT NULL,
    `key` char(64) NOT NULL,
    `value` char(64) NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_extension_on_block_number (block_number),
    KEY index_extension_on_lock_hash_crc (lock_hash_crc)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE IF NOT EXISTS extension_kv_pair_versions (
    id bigint NOT NULL AUTO_INCREMENT,
    old_block_number bigint unsigned NOT NULL,
    block_number bigint unsigned NOT NULL,
    `key` char(64) NOT NULL,
    `value` char(64) NOT NULL,
    old_value char(64) NOT NULL,
    lock_hash char(64) NOT NULL,
    action_type tinyint unsigned NOT NULL,
    tx_index int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_extension_versions_on_block_number (block_number, action_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;