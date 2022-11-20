ALTER TABLE joy_id_infos
    MODIFY credential_id varchar(100) DEFAULT '' NOT NULL;
ALTER TABLE joy_id_info_versions
    MODIFY credential_id varchar(100) DEFAULT '' NOT NULL;
ALTER TABLE sub_key_infos
    MODIFY credential_id varchar(100) DEFAULT '' NOT NULL;
