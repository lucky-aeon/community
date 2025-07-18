-- 评论表情回复系统数据库迁移文件
-- 创建时间: 2025-01-17

-- 创建表情回复表
CREATE TABLE `comment_reactions` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `comment_id` int(11) NOT NULL COMMENT '评论ID',
    `user_id` int(11) NOT NULL COMMENT '用户ID',
    `reaction_type` varchar(50) NOT NULL COMMENT '表情类型',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_comment_user_reaction` (`comment_id`, `user_id`, `reaction_type`, `deleted_at`) COMMENT '同一用户对同一评论的同种表情回复唯一',
    KEY `idx_comment_id` (`comment_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_reaction_type` (`reaction_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='评论表情回复表';

-- 创建表情类型配置表
CREATE TABLE `expression_types` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `code` varchar(50) NOT NULL COMMENT '表情代码',
    `name` varchar(100) NOT NULL COMMENT '表情名称',
    `image_path` varchar(255) NOT NULL COMMENT '表情图片路径',
    `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
    `is_active` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否启用',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_expression_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='表情类型配置表';

-- 插入默认表情数据
INSERT INTO `expression_types` (`code`, `name`, `image_path`, `sort_order`) VALUES
('good', '点赞', '/expression/good.png', 1),
('effort', '努力加油', '/expression/effort.png', 2),
('deny', '否定', '/expression/deny.png', 3),
('cry', '哭泣', '/expression/cry.png', 4),
('sleepy', '困倦', '/expression/sleepy.png', 5),
('happy', '大笑开心', '/expression/happy.png', 6),
('comfort', '安慰', '/expression/comfort.png', 7),
('shy', '害羞', '/expression/shy.png', 8),
('worship', '崇拜', '/expression/worship.png', 9),
('surprised', '惊讶', '/expression/surprised.png', 10),
('console', '明显安慰', '/expression/console.png', 11),
('sweat', '流汗', '/expression/sweat.png', 12),
('angry', '生气', '/expression/angry.png', 13),
('confused', '疑惑', '/expression/confused.png', 14),
('pray', '祈求可怜', '/expression/pray.png', 15),
('agree', '赞成', '/expression/agree.png', 16);