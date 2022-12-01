ALTER TABLE joy_id_infos
    ADD nickname varchar(255) NOT NULL;

ALTER TABLE joy_id_info_versions
    ADD old_nickname varchar(255) NOT NULL,
    ADD nickname varchar(255) NOT NULL;
