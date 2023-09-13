ALTER TABLE sub_key_infos ADD COLUMN `derivation_c_id` varchar(1500) NOT NULL DEFAULT '' AFTER `device_type`,
                          ADD COLUMN `derivation_commitment` varchar(64) NOT NULL DEFAULT '' AFTER `derivation_c_id`;

ALTER TABLE sub_key_info_versions ADD COLUMN `derivation_c_id` varchar(1500) NOT NULL DEFAULT '' AFTER `device_type`,
                                  ADD COLUMN `old_derivation_c_id` varchar(1500) NOT NULL DEFAULT '' AFTER `derivation_c_id`,
                                  ADD COLUMN `derivation_commitment` varchar(64) NOT NULL DEFAULT '' AFTER `old_derivation_c_id`,
                                  ADD COLUMN `old_derivation_commitment` varchar(64) NOT NULL DEFAULT '' AFTER `derivation_commitment`;
