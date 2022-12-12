ALTER TABLE sub_key_infos ADD CONSTRAINT uc_sub_key_infos_on_pub_key UNIQUE (pub_key);

ALTER TABLE sub_key_infos DROP COLUMN front_end;
