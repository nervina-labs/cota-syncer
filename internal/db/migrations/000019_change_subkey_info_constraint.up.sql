ALTER TABLE sub_key_infos DROP CONSTRAINT uc_sub_key_infos_on_pub_key;
ALTER TABLE sub_key_infos ADD CONSTRAINT uc_sub_key_infos_on_lock_hash_pub_key UNIQUE (lock_hash, pub_key);
