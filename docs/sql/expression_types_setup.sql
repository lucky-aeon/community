-- 表情类型配置表创建和数据插入
-- 如果表不存在，则创建

CREATE TABLE IF NOT EXISTS `expression_types` (
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

-- 插入默认表情数据（使用 INSERT IGNORE 避免重复插入）
INSERT IGNORE INTO `expression_types` (`code`, `name`, `image_path`, `sort_order`) VALUES
('thumbs-up', '点赞', '/expression/thumbs-up.png', 1),
('heart', '喜欢', '/expression/heart.png', 2),
('laugh', '大笑', '/expression/laugh.png', 3),
('wow', '惊讶', '/expression/wow.png', 4),
('thumbs-down', '不喜欢', '/expression/thumbs-down.png', 5),
('good', '点赞', '/expression/good.png', 6),
('effort', '努力加油', '/expression/effort.png', 7),
('deny', '否定', '/expression/deny.png', 8),
('cry', '哭泣', '/expression/cry.png', 9),
('sleepy', '困倦', '/expression/sleepy.png', 10),
('happy', '大笑开心', '/expression/happy.png', 11),
('comfort', '安慰', '/expression/comfort.png', 12),
('shy', '害羞', '/expression/shy.png', 13),
('worship', '崇拜', '/expression/worship.png', 14),
('surprised', '惊讶', '/expression/surprised.png', 15),
('console', '明显安慰', '/expression/console.png', 16),
('sweat', '流汗', '/expression/sweat.png', 17),
('angry', '生气', '/expression/angry.png', 18),
('confused', '疑惑', '/expression/confused.png', 19),
('pray', '祈求可怜', '/expression/pray.png', 20),
('agree', '赞成', '/expression/agree.png', 21);