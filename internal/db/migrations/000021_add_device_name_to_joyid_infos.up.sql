ALTER TABLE joy_id_info_versions
    ADD old_device_name varchar(255) NOT NULL AFTER front_end,
    ADD device_name varchar(255) NOT NULL AFTER old_device_name;

ALTER TABLE joy_id_infos
    ADD device_name varchar(255) NOT NULL AFTER front_end;

ALTER TABLE sub_key_infos
    ADD device_name varchar(255) NOT NULL AFTER front_end;


