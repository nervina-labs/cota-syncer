ALTER TABLE withdraw_cota_nft_kv_pairs
    ADD tx_hash char(64) NOT NULL AFTER out_point_crc;