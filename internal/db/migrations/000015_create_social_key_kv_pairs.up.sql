CREATE TABLE IF NOT EXISTS social_kv_pairs (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint UNSIGNED NOT NULL,
    lock_hash char(64) NOT NULL,
    lock_hash_crc int unsigned NOT NULL,
    recovery_mode tinyint NOT NULL,
    must tinyint NOT NULL,
    total tinyint NOT NULL,
    signers text NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT uc_social_on_lock_hash UNIQUE (lock_hash)
    KEY index_social_on_lock_hash_crc (lock_hash_crc),
    KEY index_social_on_block_number (block_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE IF NOT EXISTS social_kv_pair_versions (
    id bigint NOT NULL AUTO_INCREMENT,
    old_block_number bigint UNSIGNED NOT NULL,
    block_number bigint UNSIGNED NOT NULL,
    lock_hash char(64) NOT NULL,
    old_recovery_mode tinyint NOT NULL,
    recovery_mode tinyint NOT NULL,
    old_must tinyint NOT NULL,
    must tinyint NOT NULL,
    old_total tinyint NOT NULL,
    total tinyint NOT NULL,
    old_signers text NOT NULL,
    signers text NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_social_on_block_number (block_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
