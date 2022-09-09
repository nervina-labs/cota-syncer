ALTER TABLE register_cota_kv_pairs
    ADD cota_cell_id bigint unsigned NOT NULL DEFAULT 0 AFTER lock_hash,
