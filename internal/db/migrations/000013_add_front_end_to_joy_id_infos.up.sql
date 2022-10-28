ALTER TABLE joy_id_info_versions
    ADD old_front_end varchar(255) NOT NULL AFTER alg,
    ADD front_end varchar(255) NOT NULL AFTER old_front_end;

ALTER TABLE joy_id_infos
    ADD front_end varchar(255) NOT NULL AFTER alg;
