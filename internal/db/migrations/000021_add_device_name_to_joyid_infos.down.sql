ALTER TABLE joy_id_info_versions 
    DROP COLUMN device_name,
    DROP COLUMN old_device_name;

ALTER TABLE joy_id_infos DROP COLUMN device_name;

ALTER TABLE sub_key_infos DROP COLUMN device_name;

DROP TABLE IF EXISTS sub_key_info_versions;

