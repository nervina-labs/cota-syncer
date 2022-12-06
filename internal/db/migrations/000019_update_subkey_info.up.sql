ALTER TABLE sub_key_infos DROP CONSTRAINT uc_sub_key_infos_on_pub_key;

ALTER TABLE sub_key_infos
    ADD front_end varchar(255) NOT NULL AFTER alg;
