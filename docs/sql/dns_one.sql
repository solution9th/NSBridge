-- MySQL dump 10.13  Distrib 5.7.25, for macos10.14 (x86_64)
--
-- Host: localhost    Database: dns_one
-- ------------------------------------------------------
-- Server version	5.7.25

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `dns_one`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `dns_one` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `dns_one`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for auth
-- ----------------------------
DROP TABLE IF EXISTS `auth`;
CREATE TABLE `auth` (
                      `id` int(11) NOT NULL AUTO_INCREMENT,
                      `domain_key` char(36) NOT NULL,
                      `domain_secret` char(32) NOT NULL,
                      `remark` varchar(255) NOT NULL,
                      `disable` tinyint(4) NOT NULL,
                      `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                      `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for dns_domain
-- ----------------------------
DROP TABLE IF EXISTS `dns_domain`;
CREATE TABLE `dns_domain` (
                            `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '唯一id',
                            `fone_domain_id` int(11) DEFAULT NULL COMMENT '和f1关联',
                            `domain_key` char(50) NOT NULL COMMENT '创建这个域名的key',
                            `domain` varchar(255) NOT NULL COMMENT '域名',
                            `name_server` varchar(255) NOT NULL COMMENT 'ns 服务器',
                            `soa_email` varchar(50) NOT NULL,
                            `remark` varchar(50) NOT NULL COMMENT '备注',
                            `is_take_over` tinyint(4) NOT NULL COMMENT '是否接管',
                            `is_open_key` tinyint(4) NOT NULL COMMENT 'key 是否停用',
                            `record_key` char(33) NOT NULL COMMENT '操作该域名的key',
                            `record_secret` char(33) NOT NULL COMMENT 'secret ',
                            `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            PRIMARY KEY (`id`) USING BTREE,
                            UNIQUE KEY `domain_index` (`fone_domain_id`) USING HASH,
                            KEY `user_index` (`domain_key`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for dns_record
-- ----------------------------
DROP TABLE IF EXISTS `dns_record`;
CREATE TABLE `dns_record` (
                            `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
                            `domain_id` int(11) NOT NULL COMMENT '域名的id',
                            `fone_domain_id` int(11) DEFAULT NULL COMMENT 'foneid',
                            `fone_record_id` int(11) DEFAULT NULL,
                            `sub_domain` varchar(63) NOT NULL COMMENT '子域名',
                            `record_type` enum('A','AAAA','CNAME','TXT','MX','NS') NOT NULL COMMENT '记录类型',
                            `value` varchar(255) NOT NULL,
                            `line_id` int(11) NOT NULL COMMENT '线路id',
                            `ttl` int(11) NOT NULL,
                            `unit` enum('sec','min','hour','day') NOT NULL,
                            `priority` tinyint(11) unsigned NOT NULL DEFAULT '1' COMMENT '优先级, 默认是1',
                            `disable` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用, 默认是0',
                            `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for log
-- ----------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log` (
                     `log_id` int(11) NOT NULL,
                     `user_id` char(20) NOT NULL,
                     `content` varchar(255) NOT NULL,
                     `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                     PRIMARY KEY (`log_id`),
                     KEY `user_index` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                      `user_id` char(20) NOT NULL COMMENT 'i1用户名字的拼音，唯一',
                      `username` varchar(10) NOT NULL,
                      `email` varchar(50) NOT NULL,
                      `mobile` int(11) NOT NULL,
                      `login_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                      `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      PRIMARY KEY (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;

