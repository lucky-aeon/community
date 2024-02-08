alter table types
    add flag_name varchar(10) not null comment '唯一标识名';

alter table types
    add constraint types_pk_flag_name
        unique (flag_name);

alter table comments
    add from_user_id int(11) not null;
    add to_user_id int(11) not null;

-- ----------------------------
-- Table structure for member
-- ----------------------------
DROP TABLE IF EXISTS `member`;
CREATE TABLE `member` (
                          `id` int(11) NOT NULL AUTO_INCREMENT,
                          `name` varchar(255) DEFAULT NULL,
                          `desc` varchar(255) DEFAULT NULL,
                          `created_at` datetime DEFAULT NULL,
                          `updated_at` datetime DEFAULT NULL,
                          PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

alter table invite_codes
    add member_id int(11) not null;
