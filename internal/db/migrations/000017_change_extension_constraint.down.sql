ALTER TABLE extension_kv_pairs ADD CONSTRAINT uc_extension_pairs_on_key UNIQUE (`key`);
ALTER TABLE extension_kv_pairs DROP CONSTRAINT uc_extension_pairs_on_key_and_lock_hash;
