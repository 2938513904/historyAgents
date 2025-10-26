-- AI多智能体系统数据库初始化脚本

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `ai_agents` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE `ai_agents`;

-- 创建用户表（如果需要）
-- CREATE TABLE IF NOT EXISTS `users` (
--     `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
--     `username` varchar(50) NOT NULL UNIQUE,
--     `email` varchar(100) NOT NULL UNIQUE,
--     `password_hash` varchar(255) NOT NULL,
--     `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
--     `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--     PRIMARY KEY (`id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 设置权限
GRANT ALL PRIVILEGES ON `ai_agents`.* TO 'aiuser'@'%';
FLUSH PRIVILEGES;

SET FOREIGN_KEY_CHECKS = 1;