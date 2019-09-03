CREATE TABLE IF NOT EXISTS `dummy_bulk_data` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) DEFAULT NULL,
  `email` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='table for bulk insert tester';


TRUNCATE TABLE dummy_bulk_data;

