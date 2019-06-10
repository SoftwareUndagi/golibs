
CREATE TABLE `sec_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id dari group',
  `uuid` varchar(36) comment 'data UUID, user visible data', 
  `group_code` varchar(32) NOT NULL DEFAULT '' COMMENT 'kode group',
  `group_name` varchar(128) NOT NULL DEFAULT '' COMMENT 'nama group',
  `group_remark` varchar(1024) DEFAULT NULL COMMENT 'catatan dari group', 
  `usage_count` int(5) DEFAULT NULL comment 'usage count, to ease editor to enable del or not',
  `creator_name` varchar(256) DEFAULT NULL COMMENT 'audit trail, creator of data',
  `createdAt` datetime DEFAULT NULL COMMENT 'audit trail, time data created',
  `creator_ip_address` varchar(39) DEFAULT NULL COMMENT 'audit trail, ip address of data creator',
  `modified_by` varchar(256) DEFAULT NULL COMMENT 'audit trail, last user that modified data',
  `modified_by_ip` varchar(39) DEFAULT NULL COMMENT 'audit trail, IP address from where user modify data',
  `updatedAt` datetime DEFAULT NULL COMMENT 'audit trail, last modification time',
  `deletedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_group_code` (`group_code`), 
  UNIQUE KEY `unq_group_uuid` (`uuid`)
) ENGINE=InnoDB;

CREATE TABLE `sec_group_assignment` (
  `id` int(11)  NOT NULL AUTO_INCREMENT COMMENT 'id data(surrogate key)',
  `uuid` varchar(36) comment 'data UUID, user visible data', 
  `group_id` int(11)  NOT NULL COMMENT 'id dari group',
  `user_id` int(11)  NOT NULL COMMENT 'id dari user',
  `usage_count` int(11) DEFAULT NULL,
  `creator_name` varchar(256) DEFAULT NULL COMMENT 'audit trail, creator of data',
  `createdAt` datetime DEFAULT NULL COMMENT 'audit trail, time data created',
  `creator_ip_address` varchar(39) DEFAULT NULL COMMENT 'audit trail, ip address of data creator',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNQ_GROUPID_AND_USER_ID_CONSTRAINT` (`group_id`,`user_id`),
  UNIQUE KEY `unq_group_assignment_uuid` (`uuid`), 
  KEY `fk_group_assignment_to_user` (`user_id`),
  CONSTRAINT `pk_group_ass_group` FOREIGN KEY (`group_id`) REFERENCES `sec_group` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `pk_group_ass_user` FOREIGN KEY (`user_id`) REFERENCES `sec_user` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB ;


CREATE TABLE `sec_menu` (
  `id`  int(11) NOT NULL AUTO_INCREMENT COMMENT 'id dari menu',
  `uuid` varchar(36) comment 'data UUID, user visible data', 
  `menu_code` varchar(255) NOT NULL DEFAULT '' COMMENT 'kode menu',
  `parent_id`  int(11) DEFAULT NULL COMMENT 'id induk dari menu',
  `menu_label` varchar(255) NOT NULL DEFAULT '' COMMENT 'label dari menu',
  `menu_tree_code` varchar(255) DEFAULT NULL COMMENT 'tree code dari menu',
  `additional_param` varchar(256) DEFAULT NULL,
  `order_no` int(5) NOT NULL COMMENT 'urutan data',
  `i18n_key` varchar(128) DEFAULT NULL COMMENT 'key internalization',
  `route_path` varchar(256) DEFAULT NULL COMMENT 'handler menu. ini nama js function untuk handler task',
  `data_status` enum('A','D') NOT NULL DEFAULT 'A' COMMENT 'status data',
  `tree_level_position` int(11) DEFAULT NULL COMMENT 'level menu. pada level berapa data berada',
  `creator_name` varchar(256) DEFAULT NULL COMMENT 'audit trail, creator of data',
  `createdAt` datetime DEFAULT NULL COMMENT 'audit trail, time data created',
  `creator_ip_address` varchar(39) DEFAULT NULL COMMENT 'audit trail, ip address of data creator',
  `modified_by` varchar(256) DEFAULT NULL COMMENT 'audit trail, last user that modified data',
  `modified_by_ip` varchar(39) DEFAULT NULL COMMENT 'audit trail, IP address from where user modify data',
  `updatedAt` datetime DEFAULT NULL COMMENT 'audit trail, last modification time',
  `deletedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_menu_uuid` (`uuid`), 
  UNIQUE KEY `unq_menu_code` (`menu_code`)
) ENGINE=InnoDB ;


CREATE TABLE `sec_menu_assignment` (
  `id`  int(11) NOT NULL AUTO_INCREMENT COMMENT 'primary key data',
  `uuid` varchar(36) comment 'data UUID, user visible data', 
  `menu_id`  int(11) NOT NULL COMMENT 'id dari menu',
  `group_id`  int(11) NOT NULL COMMENT 'id group',
  `is_allow_create` enum('Y','N') DEFAULT 'Y' COMMENT 'flag : ijinkan new data',
  `is_allow_edit` enum('Y','N') DEFAULT 'Y' COMMENT 'flag : ijinkan edit data',
  `is_allow_erase` enum('Y','N') DEFAULT 'Y' COMMENT 'flag : ijinkan hapus data',
  `creator_name` varchar(256) DEFAULT NULL COMMENT 'audit trail, creator of data',
  `createdAt` datetime DEFAULT NULL COMMENT 'audit trail, time data created',
  `creator_ip_address` varchar(39) DEFAULT NULL COMMENT 'audit trail, ip address of data creator',
  `modified_by` varchar(256) DEFAULT NULL COMMENT 'audit trail, last user that modified data',
  `modified_by_ip` varchar(39) DEFAULT NULL COMMENT 'audit trail, IP address from where user modify data',
  `updatedAt` datetime DEFAULT NULL COMMENT 'audit trail, last modification time',
  `deletedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_menu_assignment_uuid` (`uuid`), 
  UNIQUE KEY `UNQ_MENU_ID_AND_GROUP_ID_CONSTRAINT` (`group_id`,`menu_id`),
  KEY `FK_igtwretwafo6y2tuhhmwna7yi` (`menu_id`),
  CONSTRAINT `pk_menu_ass_group` FOREIGN KEY (`group_id`) REFERENCES `sec_group` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `pk_menu_ass_menu` FOREIGN KEY (`menu_id`) REFERENCES `sec_menu` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB ;


CREATE TABLE `sec_user` (
  `id`  int(11) NOT NULL AUTO_INCREMENT COMMENT 'id dari data',
  `uuid` varchar(36) comment 'data UUID, user visible data', 
  `real_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'nama lengkap user',
  `user_name` varchar(256) NOT NULL DEFAULT '' COMMENT 'username yang di pergunakan untuk login',
  `data_status` varchar(1) DEFAULT NULL COMMENT 'status dari data',
  `avatar_url` varchar(512) DEFAULT NULL COMMENT 'url avatar user',
  `password` varchar(64) DEFAULT NULL COMMENT 'password dari user',
  `default_branch_code` varchar(16) DEFAULT NULL COMMENT 'branch dari user. default ngarah ke mana',
  `email` varchar(128) NOT NULL DEFAULT '' COMMENT 'email dari user',
  `emp_no` varchar(255) DEFAULT NULL COMMENT 'no karyawan',
  `expired_date` date DEFAULT NULL COMMENT 'kapan expiry dari data',
  `failed_login_attemps` int(11) DEFAULT NULL COMMENT 'count login failed',
  `locale` varchar(6) DEFAULT NULL COMMENT 'locale user',
  `is_locked` varchar(1) DEFAULT NULL COMMENT 'flag user locked',
  `phone1` varchar(20) DEFAULT NULL COMMENT 'no HP 1',
  `phone2` varchar(20) DEFAULT NULL COMMENT 'no HP 2',
  `remark` varchar(255) DEFAULT NULL COMMENT 'catatan dari user',
  `usage_count` int(5) not null DEFAULT 0  comment 'usage count. to ease editor for delete data. usagecount = 0 mean that data allowed to be deleted',
  `creator_name` varchar(256) DEFAULT NULL COMMENT 'audit trail, creator of data',
  `createdAt` datetime DEFAULT NULL COMMENT 'audit trail, time data created',
  `creator_ip_address` varchar(39) DEFAULT NULL COMMENT 'audit trail, ip address of data creator',
  `modified_by` varchar(256) DEFAULT NULL COMMENT 'audit trail, last user that modified data',
  `modified_by_ip` varchar(39) DEFAULT NULL COMMENT 'audit trail, IP address from where user modify data',
  `updatedAt` datetime DEFAULT NULL COMMENT 'audit trail, last modification time',
  
  `deletedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unq_user_uuid` (`uuid`), 
  UNIQUE KEY `UNQ_USERNAME_UNQ` (`user_name`),
  UNIQUE KEY `UNQ_USER_EMAIL` (`email`),
  KEY `IDX_USER_DEFAULT_BRANCH` (`default_branch_code`)
) ENGINE=InnoDB ;

create table sec_role {
    `code` varchar(16) not null comment 'ROLE code' , 
    `remark` varchar(256) not null comment 'role remark' , 
    `usage_count` int(5) not null DEFAULT 0  comment 'usage count. to ease editor for delete data. usagecount = 0 mean that data allowed to be deleted',
    `creator_name` varchar(256) DEFAULT NULL COMMENT 'audit trail, creator of data',
    `createdAt` datetime DEFAULT NULL COMMENT 'audit trail, time data created',
    `creator_ip_address` varchar(39) DEFAULT NULL COMMENT 'audit trail, ip address of data creator',
    `modified_by` varchar(256) DEFAULT NULL COMMENT 'audit trail, last user that modified data',
    `modified_by_ip` varchar(39) DEFAULT NULL COMMENT 'audit trail, IP address from where user modify data',
    `updatedAt` datetime DEFAULT NULL COMMENT 'audit trail, last modification time',
    PRIMARY KEY (`code`)

} ENGINE=InnoDB comment 'Security role'; 

create table sec_user_role {
    `id`  int(11) NOT NULL AUTO_INCREMENT COMMENT 'id dari data',
    `uuid` varchar(36) comment 'data UUID, user visible data', 
    `user_id` int(11) not null comment 'user ID of role, refer to table sec_user', 
    'role_code' varchar(16) not null comment 'user role, refer to table : sec_role' , 
    `creator_name` varchar(256) DEFAULT NULL COMMENT 'audit trail, creator of data',
    `createdAt` datetime DEFAULT NULL COMMENT 'audit trail, time data created',
    `creator_ip_address` varchar(39) DEFAULT NULL COMMENT 'audit trail, ip address of data creator',
    PRIMARY KEY (`id`) , 
    UNIQUE KEY `unq_user_role` (`user_id`,`role_code`)

} ENGINE=InnoDB comment 'Security role to user assignment' ; 

