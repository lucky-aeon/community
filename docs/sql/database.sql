/*
Navicat MySQL Data Transfer

Source Server         : 本地数据库
Source Server Version : 50536
Source Host           : 127.0.0.1:3306
Source Database       : community

Target Server Type    : MYSQL
Target Server Version : 50536
File Encoding         : 65001

Date: 2024-01-27 01:02:29
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for articles
-- ----------------------------
DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles` (
                            `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'primarykey',
                            `title` varchar(50) DEFAULT NULL,
                            `desc` longtext,
                            `user_id` int(11) NOT NULL,
                            `state` int(11) DEFAULT '0',
                            `like` int(11) DEFAULT '0',
                            `type` int(11) DEFAULT NULL,
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='it is issue or answer';

-- ----------------------------
-- Records of articles
-- ----------------------------
INSERT INTO `articles` VALUES ('1', 'Java', '内容', '3', '0', '0', '0', '2024-01-21 23:02:14', '2024-01-21 23:02:17', null);

-- ----------------------------
-- Table structure for article_relations
-- ----------------------------
DROP TABLE IF EXISTS `article_relations`;
CREATE TABLE `article_relations` (
                                     `id` int(11) NOT NULL,
                                     `parent_id` int(11) NOT NULL,
                                     `article_id` int(11) NOT NULL,
                                     `created_at` datetime DEFAULT NULL,
                                     `updated_at` datetime DEFAULT NULL,
                                     `deleted_at` datetime DEFAULT NULL,
                                     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of article_relations
-- ----------------------------

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
                            `id` int(11) NOT NULL AUTO_INCREMENT,
                            `parent_id` int(11) DEFAULT NULL,
                            `content` longtext NOT NULL,
                            `from_user_id` int(11) NOT NULL,
                            `to_user_id` int(11) NOT NULL,
                            `business_id` int(11) NOT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of comments
-- ----------------------------

-- ----------------------------
-- Table structure for files
-- ----------------------------
DROP TABLE IF EXISTS `files`;
CREATE TABLE `files` (
                         `id` int(11) NOT NULL AUTO_INCREMENT,
                         `file_key` varchar(255) NOT NULL,
                         `size` bigint(20) DEFAULT NULL,
                         `format` varchar(255) DEFAULT NULL,
                         `user_id` int(11) DEFAULT NULL,
                         `business_id` int(11) DEFAULT NULL,
                         `tenant_id` int(11) DEFAULT NULL,
                         `created_at` datetime DEFAULT NULL,
                         `updated_at` datetime DEFAULT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of files
-- ----------------------------
INSERT INTO `files` VALUES ('1', 'sadasd', '500', '', '0', '0', '0', '2024-01-17 23:37:31', '2024-01-17 23:37:31');
INSERT INTO `files` VALUES ('2', 'sadasd', '500', '', '0', '0', '0', '2024-01-20 17:50:31', '2024-01-20 17:50:31');
INSERT INTO `files` VALUES ('3', '3/2fdd675d-ab87-4259-8c44-52451500ddc0', '89031', 'image/png', '3', '1', '0', '2024-01-21 18:59:01', '2024-01-21 18:59:01');
INSERT INTO `files` VALUES ('4', '3/4356909a-14c1-4c55-a2c9-055096698444', '47077', 'image/png', '3', '1', '0', '2024-01-21 19:05:33', '2024-01-21 19:05:33');

-- ----------------------------
-- Table structure for invite_codes
-- ----------------------------
DROP TABLE IF EXISTS `invite_codes`;
CREATE TABLE `invite_codes` (
                                `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'primarykey',
                                `code` varchar(20) NOT NULL,
                                `member_Id` int(11) NOT NULL,
                                `state` tinyint(1) NOT NULL,
                                `created_at` datetime DEFAULT NULL,
                                `updated_at` datetime DEFAULT NULL,
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='Cache invite code for register';

-- ----------------------------
-- Records of invite_codes
-- ----------------------------
INSERT INTO `invite_codes` VALUES ('2', '2', '1', null, null);
INSERT INTO `invite_codes` VALUES ('3', '222', '1', '2024-01-15 22:04:17', '2024-01-15 22:05:03');
INSERT INTO `invite_codes` VALUES ('4', '123', '0', null, null);

-- ----------------------------
-- Table structure for types
-- ----------------------------
DROP TABLE IF EXISTS `types`;
create table types
(
    id            int auto_increment
        primary key,
    parent_id     int          null,
    title         varchar(255) null,
    `desc`        varchar(255) null,
    state         tinyint      null,
    sort          int          null,
    article_state varchar(255) null,
    created_at    datetime     null,
    updated_at    datetime     null,
    deleted_at    datetime     null,
    flag_name     varchar(10)  not null comment '唯一标识名',
    constraint types_pk_flag_name
        unique (flag_name)
);

-- ----------------------------
-- Records of types
-- ----------------------------
INSERT INTO luckyaeon_community.types (id, parent_id, title, `desc`, state, sort, article_state, created_at, updated_at, deleted_at, flag_name) VALUES (1, 0, '文章', 'sefsffff', null, null, null, null, '2024-01-29 23:24:02', null, 'awdad');
INSERT INTO luckyaeon_community.types (id, parent_id, title, `desc`, state, sort, article_state, created_at, updated_at, deleted_at, flag_name) VALUES (2, 1, 'Java八股', 'sefseff', null, null, null, null, '2024-01-30 23:02:43', null, 'nhv');
-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
                         `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'primarykey',
                         `name` varchar(20) DEFAULT NULL,
                         `account` varchar(20) DEFAULT NULL,
                         `password` varchar(255) DEFAULT NULL,
                         `invite_code` int(8) DEFAULT NULL COMMENT 'The invitation code in "invite_code" table used when registering',
                         `desc` longtext,
                         `avatar` varchar(255) DEFAULT NULL,
                         `created_at` datetime DEFAULT NULL,
                         `updated_at` datetime DEFAULT NULL,
                         `deleted_at` datetime DEFAULT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='community user';

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES ('1', '1', '1', '1', '1', null, 'https://t7.baidu.com/it/u=1595072465,3644073269&fm=193&f=GIF', '2024-01-12 23:36:52', '2024-01-12 23:36:52', null);
INSERT INTO `users` VALUES ('2', '1', '2', '1', '2', null, 'https://t7.baidu.com/it/u=1595072465,3644073269&fm=193&f=GIF', '2024-01-12 23:37:26', '2024-01-12 23:37:26', null);
INSERT INTO `users` VALUES ('3', 'xhy', 'xhy', '123', '123', '1234', '3/4356909a-14c1-4c55-a2c9-055096698444', '2024-01-21 12:14:39', '2024-01-21 20:17:56', null);
INSERT INTO `users` VALUES ('5', 'xhy', 'xhy', '123', '123', null, 'https://t7.baidu.com/it/u=1595072465,3644073269&fm=193&f=GIF', '2024-01-21 12:19:14', '2024-01-21 12:19:14', null);
INSERT INTO `users` VALUES ('6', 'xxxx', 'xxxx', '123', '123', null, 'https://t7.baidu.com/it/u=1595072465,3644073269&fm=193&f=GIF', '2024-01-21 12:22:34', '2024-01-21 12:22:34', null);


-- ----------------------------
-- Table structure for message_templates
-- ----------------------------
DROP TABLE IF EXISTS `message_templates`;
CREATE TABLE `message_templates` (
                                     `id` int(11) NOT NULL AUTO_INCREMENT,
                                     `content` longtext NOT NULL,
                                     `event_id` int(11) DEFAULT NULL,
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
                                  add state tinyint(1) DEFAULT '0',
                                    add type tinyint(1) DEFAULT '1',
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


-- ----------------------------
-- Table structure for subscription
-- ----------------------------
DROP TABLE IF EXISTS `subscriptions`;
CREATE TABLE `subscription` (
                                `id` int(11) NOT NULL AUTO_INCREMENT,
                                `user_id` int(11) NOT NULL,
                                `event_id` int(11) NOT NULL,
                                `businessId` int(11) NOT NULL,
                                `created_at` datetime DEFAULT NULL,
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

create table article_tags
(
    id          int auto_increment
        primary key,
    tag_name    varchar(20)  null,
    description varchar(255) null,
    user_id     int          not null comment '操作人',
    created_at  datetime     null,
    updated_at  datetime     null,
    deleted_at  datetime     null,
    constraint article_tags_tag_name_uindex
        unique (tag_name)
)
    comment '文章标签';

create table article_tag_user_relations
(
    user_id int not null comment '用户',
    tag_id  int not null comment '标签'
);

create unique index article_tag_user_relations_user_id_tag_id_uindex
    on article_tag_user_relations (user_id, tag_id);

