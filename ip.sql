/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50173
Source Host           : localhost:3306
Source Database       : seo

Target Server Type    : MYSQL
Target Server Version : 50173
File Encoding         : 65001

Date: 2017-08-26 11:05:32
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `ip`
-- ----------------------------
DROP TABLE IF EXISTS `ip`;
CREATE TABLE `ip` (
  `time` int(20) NOT NULL COMMENT '请求时间',
  `ip` varchar(100) NOT NULL COMMENT '请求IP',
  `www_host` varchar(100) NOT NULL DEFAULT '' COMMENT '请求域名',
  `www_url` varchar(100) DEFAULT NULL COMMENT '请求路径',
  `Referer` varchar(100) DEFAULT NULL COMMENT '来路',
  `Method` varchar(100) DEFAULT NULL COMMENT '请求方式',
  `User_Agent` varchar(200) DEFAULT NULL COMMENT '请求头',
  `cs` varchar(100) DEFAULT NULL COMMENT '请求次数',
  PRIMARY KEY (`time`,`ip`,`www_host`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of ip
-- ----------------------------
INSERT INTO `ip` VALUES ('1503716238', '127.0.0.1', '127.0.0.1:1000', '/11111222222/as/df/asd/f/as/df/as/df/as/df/as/d', '', 'GET', 'Mozilla_5.0 (Windows NT 6.1; WOW64) AppleWebKit_537.36 (KHTML_ like Gecko) Chrome_56.0.2924.87 Safari_537.36', '2');
INSERT INTO `ip` VALUES ('1503716241', '127.0.0.1', '127.0.0.1:1000', '/11111222222/as/df/asd/f/as/df/as/df/as/df/as/d', '', 'GET', 'Mozilla_5.0 (Windows NT 6.1; WOW64) AppleWebKit_537.36 (KHTML_ like Gecko) Chrome_56.0.2924.87 Safari_537.36', '2');
INSERT INTO `ip` VALUES ('1503716636', '127.0.0.1', '127.0.0.1:1000', '/11111222222/as/df/asd/f/as/df/as/df/as/df/as/d', '', 'GET', 'Mozilla_5.0 (Windows NT 6.1; WOW64) AppleWebKit_537.36 (KHTML_ like Gecko) Chrome_56.0.2924.87 Safari_537.36', '2');
