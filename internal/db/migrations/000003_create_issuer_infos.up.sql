CREATE TABLE IF NOT EXISTS issuer_infos (
    id bigint NOT NULL AUTO_INCREMENT,
    block_number bigint unsigned NOT NULL,
    lock_hash char(64) NOT NULL,
    version varchar(40) NOT NULL,
    `name` varchar(255) NOT NULL,
    avatar varchar(500) NOT NULL,
    description varchar(1000) NOT NULL,
    localization varchar(1000) NOT NULL,
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY index_issuer_infos_on_block_number (block_number),
    CONSTRAINT uc_issuer_infos_on_lock_hash UNIQUE (lock_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;