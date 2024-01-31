-- change invite_code table: code field type to int(8)
-- dateTime: 2024-01-11 21:07
ALTER TABLE `invite_code`
	CHANGE `code` `code` INT(8) DEFAULT NULL ;

-- ----------------------------
-- Table structure for message_templates
-- ----------------------------
DROP TABLE IF EXISTS `message_templates`;
CREATE TABLE `message_templates` (
                                     `id` int(11) NOT NULL AUTO_INCREMENT,
                                     `content` longtext NOT NULL,
                                     `created_at` datetime DEFAULT NULL,
                                     `updated_at` datetime DEFAULT NULL,
                                     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for message_states
-- ----------------------------
DROP TABLE IF EXISTS `message_states`;
CREATE TABLE `message_states` (
                                  `id` int(11) NOT NULL AUTO_INCREMENT,
                                  `content` longtext NOT NULL,
                                  `from` int(11) NOT NULL,
                                  `to` int(11) NOT NULL,
                                  `created_at` datetime DEFAULT NULL,
                                  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for message_logs
-- ----------------------------
DROP TABLE IF EXISTS `message_logs`;
CREATE TABLE `message_logs` (
                                `id` int(11) NOT NULL AUTO_INCREMENT,
                                `content` longtext NOT NULL,
                                `from` int(11) NOT NULL,
                                `to` int(11) NOT NULL,
                                `type` int(11) NOT NULL,
                                `created_at` datetime DEFAULT NULL,
                                `deleted_at` datetime DEFAULT NULL,
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
