ALTER TABLE joy_id_infos DROP COLUMN `derivation_c_id`,
                         DROP COLUMN `derivation_commitment`;

ALTER TABLE joy_id_info_versions DROP COLUMN `derivation_c_id`,
                         DROP COLUMN `derivation_commitment`;
