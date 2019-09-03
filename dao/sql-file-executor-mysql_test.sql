CREATE temporary TABLE `temp_ct_lookup_header` (
  `id` varchar(128) NOT NULL DEFAULT '' COMMENT 'id dari lookup',
  `is_cacheable` varchar(255) DEFAULT NULL COMMENT 'flag data cacheable atau tidak',
  `remark` varchar(512) DEFAULT NULL COMMENT 'catatan dalam LOV',
  `is_use_custom_sql` enum('Y','N') DEFAULT 'N' COMMENT 'flag use custom SQL atau tidak,kalau Y maka data di ambil dari sql',
  `sql_data` text COMMENT 'custom SQL. sql musti menghasilkan query yang equivalent dengan m_lookup_detail',
  `sql_version` text COMMENT 'sql statement untuk membaca versi dari data lookup. ini untuk handle cacheable item',
  `sql_data_filtered` text COMMENT 'sql dengan filter, query harus menyertakan {{codes}}',
  `code_actual_data_type` varchar(10) DEFAULT 'string' COMMENT 'actual data type dari filter. kalau filter di sertakan, kalau bukan string maka akan ada konversi tipe data',
  `version` varchar(32) DEFAULT NULL COMMENT 'versi dari LOV',
  `creator_name` varchar(128) DEFAULT NULL COMMENT 'creator data',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'waktu data di buat',
  `creator_ip_address` varchar(49) DEFAULT NULL COMMENT 'ip address data yang membuat data',
  `modified_by` varchar(128) DEFAULT NULL COMMENT 'user modifikasi terkahir data',
  `updatedAt` datetime DEFAULT NULL COMMENT 'waktu modifikasi terkahir data',
  `modified_by_ip` varchar(49) DEFAULT NULL COMMENT 'ip address modifikasi terkahir data',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='data header lookup';

CREATE temporary TABLE `temp_ct_lookup_details` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'id dari data',
  `detail_code` varchar(32) NOT NULL DEFAULT '' COMMENT 'kode lookup detail. ini untuk value',
  `lov_id` varchar(128) NOT NULL DEFAULT '' COMMENT 'id LOV header',
  `label` varchar(128) DEFAULT NULL COMMENT 'label dari lookup',
  `val_1` varchar(32) DEFAULT NULL COMMENT 'value 1',
  `val_2` varchar(32) DEFAULT NULL COMMENT 'value 2',
  `i18n_key` varchar(256) DEFAULT '' COMMENT 'key internalization dari data. di ambil dari local client',
  `seq_no` int(11) DEFAULT NULL COMMENT 'urutan data',
  `creator_name` varchar(128) DEFAULT NULL COMMENT 'creator data',
  `createdAt` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'waktu data di buat',
  `creator_ip_address` varchar(49) DEFAULT NULL COMMENT 'ip address data yang membuat data',
  `modified_by` varchar(128) DEFAULT '' COMMENT 'user modifikasi terkahir data',
  `updatedAt` datetime DEFAULT NULL COMMENT 'waktu modifikasi terkahir data',
  `modified_by_ip` varchar(49) DEFAULT NULL COMMENT 'ip address modifikasi terkahir data',
  PRIMARY KEY (`id`),
  KEY `fk_lookup_detail_to_header` (`lov_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='data detail lookup';


INSERT INTO `temp_ct_lookup_header` (`id`, `is_cacheable`, `remark`, `is_use_custom_sql`, `sql_data`, `sql_version`, `sql_data_filtered`, `code_actual_data_type`, `version`, `creator_name`, `createdAt`, `creator_ip_address`, `modified_by`, `updatedAt`, `modified_by_ip`)
VALUES
	('SIMPLE_LOOKUP', 'Y', 'lookup simple. dengan content pada lookup detail', 'N', NULL, NULL, NULL, 'string', '0.0.1', 'gede.sutarsa', '2019-06-09 20:51:22', '192.168.16.16', NULL, NULL, NULL),
	('SQL_LOOKUP', 'Y', 'Lookup with SQL', 'Y', NULL, NULL, NULL, 'string', NULL, NULL, '2019-06-09 20:54:32', NULL, NULL, NULL, NULL);


INSERT INTO `temp_ct_lookup_details` (`id`, `detail_code`, `lov_id`, `label`, `val_1`, `val_2`, `i18n_key`, `seq_no`, `creator_name`, `createdAt`, `creator_ip_address`, `modified_by`, `updatedAt`, `modified_by_ip`)
VALUES
	(2, 'L001', 'SIMPLE_LOOKUP', 'Label1', NULL, NULL, '', NULL, 'gede.sutarsa', '2019-06-09 20:53:11', NULL, '', NULL, NULL),
	(3, 'L002', 'SIMPLE_LOOKUP', 'Label 2', NULL, NULL, '', NULL, 'gede.sutarsa', '2019-06-09 20:53:11', NULL, '', NULL, NULL),
	(4, 'L003', 'SIMPLE_LOOKUP', 'Label 3', NULL, NULL, '', NULL, 'gede.sutarsa', '2019-06-09 20:53:11', NULL, '', NULL, NULL);
