ALTER TABLE joy_id_infos ADD COLUMN `derivation_c_id` varchar(1500) NOT NULL DEFAULT '' AFTER `cota_cell_id`,
                        ADD COLUMN `derivation_commitment` varchar(64) NOT NULL DEFAULT '' AFTER `derivation_c_id`;

ALTER TABLE joy_id_info_versions ADD COLUMN `derivation_c_id` varchar(1500) NOT NULL DEFAULT '' AFTER `cota_cell_id`,
                         ADD COLUMN `derivation_commitment` varchar(64) NOT NULL DEFAULT '' AFTER `derivation_c_id`;
