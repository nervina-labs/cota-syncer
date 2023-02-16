ALTER TABLE joy_id_info_versions 
    DROP COLUMN device_name,
    DROP COLUMN old_device_name,
    DROP COLUMN device_type,
    DROP COLUMN old_device_type;

ALTER TABLE joy_id_infos 
    DROP COLUMN device_name,
    DROP COLUMN device_type;

ALTER TABLE sub_key_infos
    DROP COLUMN device_name,
    DROP COLUMN device_type;

DROP TABLE IF EXISTS sub_key_info_versions;
