drop table if exists bulk_insert_study ; 

CREATE TABLE `bulk_insert_study` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`name` varchar(128) DEFAULT NULL,
	`email` varchar(256) DEFAULT NULL,
	`def_col` enum('Y','N') DEFAULT 'Y' COMMENT 'column with default column',
	`tag_username` varchar(128) DEFAULT NULL,
	PRIMARY KEY (`id`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;