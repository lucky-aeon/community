-- GitHub 登录功能数据库迁移脚本
-- 创建时间：2025-06-28

-- 创建用户 GitHub 绑定表
CREATE TABLE user_github_bindings (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    user_id INT NOT NULL COMMENT '用户ID，关联 users 表',
    github_id BIGINT NOT NULL UNIQUE COMMENT 'GitHub 用户ID',
    github_username VARCHAR(100) NOT NULL COMMENT 'GitHub 用户名',
    github_email VARCHAR(255) COMMENT 'GitHub 邮箱',
    github_avatar VARCHAR(500) COMMENT 'GitHub 头像URL',
    access_token VARCHAR(255) COMMENT 'GitHub Access Token（可选存储）',
    bound_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '绑定时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户GitHub账号绑定表';

-- 确保每个用户只能绑定一个 GitHub 账号
ALTER TABLE user_github_bindings ADD UNIQUE KEY uk_user_id (user_id);

-- 确保每个 GitHub 账号只能绑定一个用户
-- github_id 已经设置为 UNIQUE，这里不需要额外约束