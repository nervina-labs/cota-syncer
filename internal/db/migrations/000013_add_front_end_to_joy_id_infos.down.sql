ALTER TABLE joy_id_info_versions 
    DROP COLUMN front_end,
    DROP COLUMN old_front_end;

ALTER TABLE joy_id_infos DROP COLUMN front_end;
