ALTER TABLE withdraw_cota_nft_kv_pairs
    ADD tx_hash char(64) AFTER out_point_crc;

ALTER TABLE withdraw_cota_nft_kv_pairs
    ADD lock_script_id bigint NOT NULL AFTER lock_hash_crc;
