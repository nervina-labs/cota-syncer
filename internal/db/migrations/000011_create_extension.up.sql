CREATE TABLE IF NOT EXISTS extension (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    lock_hash char(64) NOT NULL,
    `key` char(64) NOT NULL,
    `value` char(64) NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_extension_on_block_number (block_number),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
