-- AI评论总结功能数据库表创建脚本
-- 执行时间：2025年
-- 说明：为评论系统添加AI智能总结功能

CREATE TABLE comment_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    business_id INT NOT NULL COMMENT '业务对象ID',
    tenant_id INT NOT NULL COMMENT '租户ID: 0=文章 1=章节 2=课程 3=分享会 4=AI日报',
    summary TEXT NOT NULL COMMENT 'AI生成的总结内容',
    comment_count INT DEFAULT 0 COMMENT '参与总结的评论数量',
    last_comment_id INT DEFAULT 0 COMMENT '最后处理的评论ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_business_tenant (business_id, tenant_id),
    INDEX idx_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论AI总结表';