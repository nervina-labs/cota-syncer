ALTER TABLE define_cota_nft_kv_pairs
    ADD lock_script_id bigint NOT NULL AFTER lock_hash_crc;
