-- 通用表情回复表
CREATE TABLE reactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    business_type INT NOT NULL COMMENT '业务类型: 0=文章, 1=评论, 2=课程, 3=分享会, 4=AI日报',
    business_id BIGINT NOT NULL COMMENT '业务ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    reaction_type VARCHAR(50) NOT NULL COMMENT '表情类型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    UNIQUE KEY uk_user_business_reaction (user_id, business_type, business_id, reaction_type),
    KEY idx_business (business_type, business_id),
    KEY idx_user (user_id),
    KEY idx_reaction_type (reaction_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='通用表情回复表';

-- 数据迁移：将现有评论表情数据迁移到新表
INSERT INTO reactions (business_type, business_id, user_id, reaction_type, created_at, updated_at, deleted_at)
SELECT 1 as business_type, comment_id, user_id, reaction_type, created_at, updated_at, deleted_at
FROM comment_reactions;

-- 验证迁移完成后删除旧表（先注释掉，待验证后再执行）
-- DROP TABLE comment_reactions;