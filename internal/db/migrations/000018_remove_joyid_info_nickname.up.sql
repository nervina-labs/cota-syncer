ALTER TABLE joy_id_infos DROP COLUMN nickname;

ALTER TABLE joy_id_info_versions 
    DROP COLUMN old_nickname,
    DROP COLUMN nickname;
