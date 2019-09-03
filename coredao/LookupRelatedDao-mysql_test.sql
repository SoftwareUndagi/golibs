

CREATE TABLE IF NOT EXISTS `sec_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id dari data',
  `soverign_auth_id` varchar(128) DEFAULT NULL COMMENT 'ID of other auth. for example firebase ID',
  `uid` varchar(36) DEFAULT NULL COMMENT 'data UUID, user visible data',
  `real_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'nama lengkap user',
  `user_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'username yang di pergunakan untuk login',
  `avatar_url` varchar(512) DEFAULT NULL COMMENT 'url avatar user',
  `user_password` varchar(512) DEFAULT NULL COMMENT 'password dari user',
  `default_branch_code` varchar(16) DEFAULT NULL COMMENT 'branch dari user. default ngarah ke mana',
  `email` varchar(512) NOT NULL DEFAULT '' COMMENT 'email dari user',
  `emp_no` varchar(255) DEFAULT NULL COMMENT 'no karyawan',
  `expired_date` date DEFAULT NULL COMMENT 'kapan expiry dari data',
  `failed_login_attemps` int(11) DEFAULT NULL COMMENT 'count login failed',
  `locale_code` varchar(6) DEFAULT NULL COMMENT 'locale user',
  `is_locked` enum('Y','N') DEFAULT NULL COMMENT 'flag user locked',
  `is_active` enum('Y','N') DEFAULT NULL COMMENT 'status dari data',
  `phone1` varchar(20) DEFAULT NULL COMMENT 'no HP 1',
  `phone2` varchar(20) DEFAULT NULL COMMENT 'no HP 2',
  `remark` varchar(255) DEFAULT NULL COMMENT 'catatan dari user',
  `usage_count` int(5) NOT NULL DEFAULT '0' COMMENT 'usage count. to ease editor for delete data. usagecount = 0 mean that data allowed to be deleted',
  `creator_name` varchar(256) DEFAULT NULL COMMENT 'audit trail, creator of data',
  `createdAt` datetime DEFAULT NULL COMMENT 'audit trail, time data created',
  `creator_ip_address` varchar(39) DEFAULT NULL COMMENT 'audit trail, ip address of data creator',
  `modified_by` varchar(256) DEFAULT NULL COMMENT 'audit trail, last user that modified data',
  `modified_by_ip` varchar(39) DEFAULT NULL COMMENT 'audit trail, IP address from where user modify data',
  `updatedAt` datetime DEFAULT NULL COMMENT 'audit trail, last modification time',
  `deletedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNQ_USERNAME_UNQ` (`user_name`),
  UNIQUE KEY `UNQ_USER_EMAIL` (`email`),
  UNIQUE KEY `unq_user_uuid` (`uid`),
  KEY `IDX_USER_DEFAULT_BRANCH` (`default_branch_code`)
) ENGINE=InnoDB AUTO_INCREMENT=156 DEFAULT CHARSET=utf8;



CREATE TABLE IF NOT EXISTS `ct_lookup_header` (
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

CREATE TABLE IF NOT EXISTS `ct_lookup_details` (
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
  KEY `fk_lookup_detail_to_header` (`lov_id`),
  CONSTRAINT `fk_lookup_detail_to_header` FOREIGN KEY (`lov_id`) REFERENCES `ct_lookup_header` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=233 DEFAULT CHARSET=utf8 COMMENT='data detail lookup';

delete from sec_user where user_name in ('dev-test-helper.gede.sutarsa@gmail.com' ,'dev-test-helper.yuli.dana@gmail.com' ) ; 
delete from ct_lookup_details where lov_id in ('SIMPLE_LOOKUP'   ) ; 
delete from ct_lookup_header where id in ('SIMPLE_LOOKUP' , 'TEST::MasterUser'  ) ; 

INSERT INTO `ct_lookup_header` (`id`, `is_cacheable`, `remark`, `is_use_custom_sql`, `sql_data`, `sql_version`, `sql_data_filtered`, `code_actual_data_type`, `version`, `creator_name`, `createdAt`, `creator_ip_address`, `modified_by`, `updatedAt`, `modified_by_ip`)
VALUES
	('SIMPLE_LOOKUP', 'Y', 'lookup simple. dengan content pada lookup detail', 'N', NULL, NULL, NULL, 'string', '0.0.1', 'gede.sutarsa', '2019-06-09 20:51:22', '192.168.16.16', NULL, NULL, NULL),
	('TEST::MasterUser', 'Y', 'Lookup with SQL(user)', 'Y', 'select \n 1 ID , \'gede.sutarsa\' CreatorName , \'127.0.0.1\' CreatorIPAddress , null CreatedAt , null UpdatedAt, null ModifiedBy, null ModifiedIPAddress , 	m.user_name DetailCode , \n	\'TEST::MasterUser\' LovID , \n	m.real_name  Label , \n	m.email   Value1 ,  \n	m.locale_code  Value2  , \n	null  I18nKey  , \n	@rownum:=ifnull(@rownum , 0 )+1 SequenceNo\nfrom \n	sec_user m \nwhere\n	m.user_name in (\'dev-test-helper.gede.sutarsa@gmail.com\' , \'dev-test-helper.yuli.dana@gmail.com\') \norder by  m.user_name', 'select \'TEST::MasterUser\' lovId,\nleft(convert(  		\ncase when  max(ifnull(createdAt,\'0001-1-1\'))> max(ifnull(updatedAt,\'0001-01-01\'))   then \n	UNIX_TIMESTAMP(max(ifnull(createdAt,\'0001-1-1\')))\nelse\n	UNIX_TIMESTAMP(max(ifnull(updatedAt,\'0001-01-01\')))\nend  , char ),10)version  from sec_user m \nwhere\n	1=1\n	and m.user_name in (\'dev-test-helper.gede.sutarsa@gmail.com\' , \'dev-test-helper.yuli.dana@gmail.com\') ', 'set @rownum = 0;\nselect \n	m.user_name detailCode , \n	\'TEST::MasterUser\' lovId , \n	m.real_name  label , \n	m.email   value1 ,  \n	m.locale_code  value2  , \n	null  i18nKey  , \n	@rownum:=ifnull(@rownum , 0 )+1 sequenceNo\nfrom \n	sec_user m \nwhere\n	1=1\n	and m.user_name in (\'dev-test-helper.gede.sutarsa@gmail.com\' , \'dev-test-helper.yuli.dana@gmail.com\') \n	and  m.user_name in ({{codes}})	\norder by  m.user_name', 'string', NULL, NULL, '2019-06-09 20:54:32', NULL, NULL, NULL, NULL);


INSERT INTO `ct_lookup_details` ( `detail_code`, `lov_id`, `label`, `val_1`, `val_2`, `i18n_key`, `seq_no`, `creator_name`, `createdAt`, `creator_ip_address`, `modified_by`, `updatedAt`, `modified_by_ip`)
VALUES
	(  'L001', 'SIMPLE_LOOKUP', 'Label1', NULL, NULL, '', NULL, 'gede.sutarsa', '2019-06-09 20:53:11', NULL, '', NULL, NULL),
	(  'L002', 'SIMPLE_LOOKUP', 'Label 2', NULL, NULL, '', NULL, 'gede.sutarsa', '2019-06-09 20:53:11', NULL, '', NULL, NULL),
	(  'L003', 'SIMPLE_LOOKUP', 'Label 3', NULL, NULL, '', NULL, 'gede.sutarsa', '2019-06-09 20:53:11', NULL, '', NULL, NULL);


INSERT INTO `sec_user` ( `soverign_auth_id`, `uid`, `real_name`, `user_name`, `avatar_url`, `user_password`, `default_branch_code`, `email`, `emp_no`, `expired_date`, `failed_login_attemps`, `locale_code`, `is_locked`, `is_active`, `phone1`, `phone2`, `remark`, `usage_count`, `creator_name`, `createdAt`, `creator_ip_address`, `modified_by`, `modified_by_ip`, `updatedAt`, `deletedAt`)
VALUES
	(  'GEDE101', 'UUIDUSR001', 'Gede Sutarsa', 'dev-test-helper.gede.sutarsa@gmail.com', NULL, NULL, NULL, 'gede.sutarsa@gmail.com', NULL, NULL, NULL, NULL, 'N', 'Y', NULL, NULL, NULL, 0, NULL, current_timestamp, NULL, NULL, NULL, NULL, NULL),
	(  'YULI001', 'UUIDUSR002', 'Yuliati Danayanti', 'dev-test-helper.yuli.dana@gmail.com', NULL, NULL, NULL, 'yuli.dana@gmail.com', NULL, NULL, NULL, NULL, 'N', 'Y', NULL, NULL, NULL, 0, NULL, current_timestamp, NULL, NULL, NULL, NULL, NULL);
