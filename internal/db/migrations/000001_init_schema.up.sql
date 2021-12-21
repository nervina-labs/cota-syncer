START TRANSACTION;

CREATE TABLE IF NOT EXISTS check_infos (
    id bigint NOT NULL AUTO_INCREMENT,
    check_type tinyint unsigned NOT NULL,
    block_number bigint unsigned NOT NULL,
    block_hash char(64) NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_check_infos_on_block_number (block_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS register_cota_kv_pairs (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    lock_hash char(64) NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_register_on_block_number (block_number),
    CONSTRAINT uc_register_on_lock_hash UNIQUE (lock_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS define_cota_nft_kv_pairs (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    cota_id char(40) NOT NULL,
    total int unsigned NOT NULL,
    issued int unsigned NOT NULL,
    configure tinyint unsigned NOT NULL,
    lock_hash char(64) NOT NULL,
    lock_hash_crc int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_define_on_block_number (block_number),
    KEY index_define_on_lock_hash_crc (lock_hash_crc),
    CONSTRAINT uc_define_on_cota_id UNIQUE (cota_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS define_cota_nft_kv_pair_versions (
    id bigint NOT NULL AUTO_INCREMENT,
    old_block_number bigint unsigned NOT NULL,
    block_number bigint unsigned NOT NULL,
    cota_id char(40) NOT NULL,
    total int unsigned NOT NULL,
    old_issued int unsigned NOT NULL,
    issued int unsigned NOT NULL,
    configure tinyint unsigned NOT NULL,
    lock_hash char(64) NOT NULL,
    action_type tinyint unsigned NOT NULL,
    tx_index int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_define_versions_on_block_number (block_number, action_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS hold_cota_nft_kv_pairs (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    cota_id char(40) NOT NULL,
    token_index int unsigned NOT NULL,
    state tinyint unsigned NOT NULL,
    configure tinyint unsigned NOT NULL,
    characteristic char(40) NOT NULL,
    lock_hash char(64) NOT NULL,
    lock_hash_crc int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_hold_on_block_number (block_number),
    KEY index_hold_on_lock_hash_crc (lock_hash_crc),
    CONSTRAINT uc_hold_on_cota_id_and_token_index UNIQUE (cota_id, token_index)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS hold_cota_nft_kv_pair_versions (
    id bigint NOT NULL AUTO_INCREMENT,
    old_block_number bigint unsigned NOT NULL,
    block_number bigint unsigned NOT NULL,
    cota_id char(40) NOT NULL,
    token_index int unsigned NOT NULL,
    old_state tinyint unsigned NOT NULL,
    state tinyint unsigned NOT NULL,
    configure tinyint unsigned NOT NULL,
    old_characteristic char(40) NOT NULL,
    characteristic char(40) NOT NULL,
    old_lock_hash char(64) NOT NULL,
    lock_hash char(64) NOT NULL,
    action_type tinyint unsigned NOT NULL,
    tx_index int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_hold_versions_on_block_number (block_number, action_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS withdraw_cota_nft_kv_pairs (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    cota_id char(40) NOT NULL,
    cota_id_crc int unsigned NOT NULL,
    token_index int unsigned NOT NULL,
    out_point char(72) NOT NULL,
    out_point_crc int unsigned NOT NULL,
    state tinyint unsigned NOT NULL,
    configure tinyint unsigned NOT NULL,
    characteristic char(40) NOT NULL,
    receiver_lock_hash char(64) NOT NULL,
    receiver_lock_hash_crc int unsigned NOT NULL,
    lock_hash char(64) NOT NULL,
    lock_hash_crc int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_withdraw_on_block_number (block_number),
    KEY index_withdraw_on_cota_id_crc_token_index (cota_id_crc, token_index),
    KEY index_withdraw_on_receiver_lock_hash_crc (receiver_lock_hash_crc),
    KEY index_withdraw_on_out_point_crc (out_point_crc),
    KEY index_withdraw_on_lock_hash_crc (lock_hash_crc)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS claimed_cota_nft_kv_pairs (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    cota_id char(40) NOT NULL,
    cota_id_crc int unsigned NOT NULL,
    token_index int unsigned NOT NULL,
    out_point char(72) NOT NULL,
    out_point_crc int unsigned NOT NULL,
    lock_hash char(64) NOT NULL,
    lock_hash_crc int unsigned NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_claimed_on_block_number (block_number),
    KEY index_claimed_on_cota_id_crc_token_index (cota_id_crc, token_index),
    KEY index_claimed_on_out_point_crc (out_point_crc),
    KEY index_claimed_on_lock_hash_crc (lock_hash_crc)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

COMMIT;