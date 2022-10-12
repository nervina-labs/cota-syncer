ALTER TABLE register_cota_kv_pairs
    ADD lock_script_id bigint NOT NULL DEFAULT 3094967296 AFTER lock_hash;
