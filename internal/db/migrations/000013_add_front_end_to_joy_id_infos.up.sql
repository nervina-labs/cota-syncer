ALTER TABLE joy_id_info_versions
    ADD front_end varchar(255) NOT NULL AFTER alg;

ALTER TABLE joy_id_infos
    ADD front_end varchar(255) NOT NULL AFTER alg;
