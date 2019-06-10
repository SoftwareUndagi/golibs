package dao

import "fmt"

var SecGroupTempTableName = "sec_group_temp"
var SecGroupAssignmentTempTableName = "sec_menu_assignment_temp"
var TempUserTempTableName = "sec_user_temp"
var TempUserGroupAssignmentTableName = "sec_group_assignment_temp"
var TempUserRoleTableName = "sec_user_role"
var TempRoleTableName = "sec_role_temp"

//TempGroupTableCreateSQL sql to create temp group table
var TempGroupTableCreateMySQL = fmt.Sprintf(`CREATE  TEMPORARY TABLE  %s (
	id int(11) NOT NULL   COMMENT 'id dari group',
	uid varchar(36) DEFAULT NULL COMMENT 'data UUID, user visible data',
	code varchar(32) NOT NULL DEFAULT '' COMMENT 'kode group',
	name varchar(128) NOT NULL DEFAULT '' COMMENT 'nama group',
	remark varchar(1024) DEFAULT NULL COMMENT 'catatan dari group' 
  )  ;`, SecGroupTempTableName)

//
var TempGroupAssignmentCreateMysql = fmt.Sprintf(`CREATE  TEMPORARY TABLE %s (
	id int(11) NOT NULL  COMMENT 'primary key data',
	uid varchar(36) DEFAULT NULL COMMENT 'data UUID, user visible data',
	menu_id int(11) NOT NULL COMMENT 'id dari menu',
	group_id int(11) NOT NULL COMMENT 'id group',
	is_allow_create enum('Y','N') DEFAULT 'Y' COMMENT 'flag : ijinkan new data',
	is_allow_edit enum('Y','N') DEFAULT 'Y' COMMENT 'flag : ijinkan edit data',
	is_allow_erase enum('Y','N') DEFAULT 'Y' COMMENT 'flag : ijinkan hapus data' 
  )   `, SecGroupAssignmentTempTableName)

var TempUserCreateMysql = fmt.Sprintf(`CREATE temporary TABLE %s (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id dari data',
	soverign_auth_id varchar(128) DEFAULT NULL COMMENT 'ID of other auth. for example firebase ID',
	uid varchar(36) DEFAULT NULL COMMENT 'data UUID, user visible data',
	real_name varchar(256) NOT NULL DEFAULT '' COMMENT 'nama lengkap user',
	user_name varchar(256) NOT NULL DEFAULT '' COMMENT 'username yang di pergunakan untuk login',
	avatar_url varchar(512) DEFAULT NULL COMMENT 'url avatar user',
	user_password varchar(512) DEFAULT NULL COMMENT 'password dari user',
	default_branch_code varchar(16) DEFAULT NULL COMMENT 'branch dari user. default ngarah ke mana',
	email varchar(512) NOT NULL DEFAULT '' COMMENT 'email dari user',
	emp_no varchar(255) DEFAULT NULL COMMENT 'no karyawan',
	expired_date date DEFAULT NULL COMMENT 'kapan expiry dari data',
	failed_login_attemps int(11) DEFAULT NULL COMMENT 'count login failed',
	locale_code varchar(6)   COMMENT 'locale user',
	is_locked enum('Y','N')   COMMENT 'flag user locked',
	is_active enum('Y','N')   COMMENT 'status dari data',
	phone1 varchar(20)   COMMENT 'no HP 1',
	phone2 varchar(20)   COMMENT 'no HP 2',
	remark varchar(255)   COMMENT 'catatan dari user',
	PRIMARY KEY (id),
	UNIQUE KEY UNQ_USERNAME_UNQ (user_name),
	UNIQUE KEY UNQ_USER_EMAIL (email),
	UNIQUE KEY unq_user_uuid (uid),
	KEY IDX_USER_DEFAULT_BRANCH (default_branch_code)
  )  `, TempUserTempTableName)

var TempUserGroupAssignmentCreateMySQL = fmt.Sprintf(`CREATE temporary TABLE %s (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id data(surrogate key)',
	uid varchar(36) DEFAULT NULL COMMENT 'data UUID, user visible data',
	group_id int(11) NOT NULL COMMENT 'id dari group',
	user_id int(11) NOT NULL COMMENT 'id dari user',
	PRIMARY KEY (id),
	UNIQUE KEY UNQ_GROUPID_AND_USER_ID_CONSTRAINT (group_id,user_id),
	UNIQUE KEY unq_group_assignment_uuid (uid),
	KEY fk_group_assignment_to_user (user_id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`, TempUserGroupAssignmentTableName)
var TempUserRoleCreateMysql = fmt.Sprintf(`CREATE temporary TABLE %s (
	id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id dari data',
	uid varchar(36) DEFAULT NULL COMMENT 'data UUID, user visible data',
	user_id int(11) NOT NULL COMMENT 'user ID of role, refer to table sec_user',
	role_code varchar(16) NOT NULL COMMENT 'user role, refer to table : sec_role',
	PRIMARY KEY (id),
	UNIQUE KEY unq_user_role (user_id,role_code)
  ) ENGINE=InnoDB  ;`, TempUserRoleTableName)
var TempRoleCreateMysql = fmt.Sprintf(`CREATE temporary TABLE %s (
	code varchar(16) NOT NULL COMMENT 'ROLE code',
	remark varchar(256) NOT NULL COMMENT 'role remark',
	PRIMARY KEY (code)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Security role';`, TempRoleTableName)
var TempGroupInsertSQL = fmt.Sprintf(`INSERT INTO %s (id, uid, code, name, remark )
VALUES
	(1, 'UUID001', 'G1', 'Grup1', 'Group test 1' ),
	(2, 'UUID002', 'G2', 'Grup2', 'Group test 2' );
`, SecGroupTempTableName)

var TempGroupAssignInsertMySQL = fmt.Sprintf(`INSERT INTO %s (id, uid, menu_id, group_id, is_allow_create, is_allow_edit, is_allow_erase )
VALUES
	(1, 'ASGN001', 1, 1, 'Y', 'Y', 'Y' ),
	(2, 'ASGN002', 2, 1, 'Y', 'Y', 'Y' ),
	(3, 'ASGN0B1', 1, 2, 'Y', 'Y', 'Y' ),
	(5, 'ASGN0B2', 2, 2, 'Y', 'Y', 'Y' );
`, SecGroupAssignmentTempTableName)

//TempUserInsertMySQL populate table temp user
var TempUserInsertMySQL = fmt.Sprintf(`INSERT INTO %s (id, soverign_auth_id, uid, real_name, user_name, avatar_url, user_password, default_branch_code, email, emp_no, expired_date, failed_login_attemps, locale_code, is_locked, is_active, phone1, phone2, remark)
VALUES
	(1, 'GEDE101', 'UUIDUSR001', 'Gede Sutarsa', 'gede.sutarsa@gmail.com', NULL, NULL, NULL, 'gede.sutarsa@gmail.com', NULL, NULL, NULL, NULL, 'N', 'Y', NULL, NULL, NULL),
	(2, 'YULI001', 'UUIDUSR002', 'Yuliati Danayanti', 'yuli.dana@gmail.com', NULL, NULL, NULL, 'yuli.dana@gmail.com', NULL, NULL, NULL, NULL, 'N', 'Y', NULL, NULL, NULL);
`, TempUserTempTableName)
var TempUserAssignmentInsertMySQL = fmt.Sprintf(`INSERT INTO %s (id, uid, group_id, user_id )
VALUES
	(1, 'UUIDG1', 1, 1),
	(2, 'UUIDG2', 2, 1),
	(3, 'UUIDG3', 3, 1),
	(6, 'UUIDG4', 1, 3),
	(8, 'UUIDG5', 2, 3);
`, TempUserGroupAssignmentTableName)
var TempRoleInsertMySQL = fmt.Sprintf(`INSERT INTO  %s (code, remark)
VALUES
	('ADMIN', 'Admin user'),
	('ROLE001', 'Role tes1'),
	('ROLE002', 'Role tes2');
`, TempRoleTableName)
var TempUserRoleInsertMySQL = fmt.Sprintf(`INSERT INTO %s (id, uid, user_id, role_code )
VALUES
	(1, 'UUIDRL001', 1, 'ADMIN' ),
	(2, 'UUIDRL002', 1, 'ROLE001' ),
	(3, 'UUIDRL003', 3, 'ROLE001' );
`, TempUserRoleTableName)
