CREATE TABLE IF NOT EXISTS token_class_audios (
    id bigint NOT NULL AUTO_INCREMENT,
    cota_id char(40) NOT NULL,
    url varchar(255) NOT NULL DEFAULT '' COMMENT 'audio url',
    name varchar(255) NOT NULL DEFAULT '' COMMENT 'name',
    idx  tinyint unsigned  NOT NULL DEFAULT 0  COMMENT 'idx',
    created_at datetime(6) NOT NULL,
    updated_at datetime(6) NOT NULL,
    PRIMARY KEY (id),
    KEY `idx_cota_id_on_class_audios` (`cota_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
