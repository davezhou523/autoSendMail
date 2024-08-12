/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80012 (8.0.12)
 Source Host           : localhost:3306
 Source Schema         : trade

 Target Server Type    : MySQL
 Target Server Version : 80012 (8.0.12)
 File Encoding         : 65001

 Date: 11/08/2024 20:42:57
*/
create  database  if not exists trade;
use trade;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for attach
-- ----------------------------
DROP TABLE IF EXISTS `attach`;

CREATE TABLE `attach`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件名',
  `file_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文件路径',
  `type` int(11) NULL DEFAULT NULL COMMENT '类型:1：邮件内容',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of attach
-- ----------------------------
INSERT INTO `attach` VALUES (1, 'content1.png', '/static/email/content1.png', 1, '2024-08-10 11:22:26', NULL);
INSERT INTO `attach` VALUES (2, 'content2-1.png', '/static/email/content2-1.png', 1, '2024-08-10 11:22:26', NULL);
INSERT INTO `attach` VALUES (3, 'content2-2.png', '/static/email/content2-2.png', 1, '2024-08-10 11:22:26', NULL);
INSERT INTO `attach` VALUES (4, 'content2-3.png', '/static/email/content2-3.png', 1, '2024-08-10 11:22:26', NULL);
INSERT INTO `attach` VALUES (5, 'content2-4.png', '/static/email/content2-4.png', 1, '2024-08-10 11:22:26', NULL);
INSERT INTO `attach` VALUES (6, 'content3.png', '/static/email/content3.png', 1, '2024-08-10 11:22:26', NULL);
INSERT INTO `attach` VALUES (7, 'content4-1.png', '/static/email/content4-1.png', 1, '2024-08-10 11:22:26', NULL);
INSERT INTO `attach` VALUES (8, 'content4-2.png', '/static/email/content4-2.png', 1, '2024-08-10 11:22:26', NULL);
INSERT INTO `attach` VALUES (9, 'content4-3.png', '/static/email/content4-3.png', 1, '2024-08-10 11:22:26', NULL);

-- ----------------------------
-- Table structure for email_content
-- ----------------------------
DROP TABLE IF EXISTS `email_content`;
CREATE TABLE `email_content`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '邮件标题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '邮件内容',
  `attach_id` json NULL COMMENT '附件id',
  `sort` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '发送邮件顺序',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sort`(`sort` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '邮件模板' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of email_content
-- ----------------------------
INSERT INTO `email_content` VALUES (1, 'Certified Quality: Trust Our ISO and CE Certified Gloves', '<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Title</title>\r\n</head>\r\n<body style=\" font-family: Arial, sans-serif;\r\n            background-color: #f5f5f5;\r\n            color: #333333;\r\n            margin: 0;\r\n            padding: 20px;\">\r\n<div style=\"font-size: 16px;\r\n            line-height: 1.6;\r\n            background-color: #ffffff;\r\n            max-width: 600px;\r\n            padding: 20px;\r\n            border-radius: 8px;\r\n            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">\r\n<p>Dear Manager,</p>\r\n<p>Hello!</p>\r\n<p>This is Una from <span style=\"font-weight: bold\">Sunwei Industrial Co., Ltd.</span>, specializing in the production of high-quality gloves, including <br>\r\n    <span style=\"font-weight: bold\">Latex gloves, Nitrile gloves, Vinyl and PE gloves.</span> Our gloves are <span style=\"font-weight: bold\">ISO and CE </span>certified with guaranteed quality,\r\n    <br>which widely used in medical, food processing, cleaning, and other industries.\r\n</p>\r\n<p>Our stable and competitive pricing can help you better predict and control your costs.</p>\r\n<p>We believe our products can meet your company\'s needs and help improve your business efficiency.</p>\r\n<p>Any interest,samples with more product details and the latest price could be sent for your reference.</p>\r\n<p>We look forward to the opportunity to collaborate with you.</p>\r\n<img src=\"cid:content1.png\" alt=\"glove image\">\r\n<p>Best regards,<br>\r\nUna Huang\r\n</p>\r\n<p></p>\r\n<p>Sunwei Industrial Co., Ltd.<br>\r\nEmail: sales2@sunweiglove.com<br>\r\nTell & Whats app:+86-13751336933<br>\r\nECONOMIC DEVELOPMENT ZONE,CAOXIAN COUNTY, SHANDONG PROVINCE, P.R. CHINA.<br>\r\n    <span style=\"color: yellowgreen\">Please consider the environment before printing this e-mail.</span><br>\r\n</p>\r\n</div>\r\n</body>\r\n</html>\r\n', '[1]', 1, '2024-08-10 11:32:54', NULL);
INSERT INTO `email_content` VALUES (2, 'Experience Exceptional Customer Support for Your Glove Orders', '<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Title</title>\r\n</head>\r\n<body style=\" font-family: Arial, sans-serif;\r\n            background-color: #f5f5f5;\r\n            color: #333333;\r\n            margin: 0;\r\n            padding: 20px;\">\r\n<div style=\"font-size: 16px;\r\n            line-height: 1.6;\r\n            background-color: #ffffff;\r\n            max-width: 600px;\r\n            padding: 20px;\r\n            border-radius: 8px;\r\n            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">\r\n<p>Dear</p>\r\n<p>Greetings!</p>\r\n <p>We are <span style=\"font-weight: bold\">Sunwei Industrial Co., Ltd</span>, specializing in the production of <span style=\"font-weight: bold\"> high-quality disposable gloves</span>. <br>\r\n    Our gloves are durable, comfortable, and suitable for various industry needs.\r\n</p>\r\n<p>We know that unreliable supply chains can cause severe disruptions. Our efficient production and logistics ensure a steady and reliable supply of high-quality gloves, even during peak demand periods.\r\n        <br>We pride ourselves on our ability to maintain a consistent supply of high-quality disposable gloves.\r\n</p>\r\n<p>We hope to become your trusted partner, ensuring you never face shortages again,please reach out to discuss how we can benefit your company.</p>\r\n<p> Looking forward to your contact.</p>\r\n<img src=\"cid:content2-1.png\" alt=\"glove image\"><img src=\"cid:content2-2.png\" alt=\"glove image\">\r\n<img src=\"cid:content2-3.png\" alt=\"glove image\"><img src=\"cid:content2-4.png\" alt=\"glove image\">\r\n<p>Best regards,<br>\r\nUna Huang\r\n</p>\r\n<p></p>\r\n<p>Sunwei Industrial Co., Ltd.<br>\r\nEmail: sales2@sunweiglove.com<br>\r\nTell & Whats app:+86-13751336933<br>\r\nECONOMIC DEVELOPMENT ZONE,CAOXIAN COUNTY, SHANDONG PROVINCE, P.R. CHINA.<br>\r\n    <span style=\"color: yellowgreen\">Please consider the environment before printing this e-mail.</span><br>\r\n</p>\r\n</div>\r\n</body>\r\n</html>\r\n', '[2, 3, 4, 5]', 2, '2024-08-10 11:32:57', NULL);
INSERT INTO `email_content` VALUES (3, 'Lock in Stable Pricing for Your Disposable Glove Needs', '<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Title</title>\r\n</head>\r\n<body style=\" font-family: Arial, sans-serif;\r\n            background-color: #f5f5f5;\r\n            color: #333333;\r\n            margin: 0;\r\n            padding: 20px;\">\r\n<div style=\"font-size: 16px;\r\n            line-height: 1.6;\r\n            background-color: #ffffff;\r\n            max-width: 600px;\r\n            padding: 20px;\r\n            border-radius: 8px;\r\n            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">\r\n<p>Dear</p>\r\n<p> I am a sales representative from <span style=\"font-weight: bold\">Sunwei Industrial Co., Ltd</span>. We understand that inconsistent quality in disposable gloves can disrupt your operations and lead to significant issues. Our gloves undergo rigorous quality control to ensure that each pair meets the highest standards.\r\n</p>\r\n<p>We would love the opportunity to provide you with reliable and high-quality gloves,as well as competitive prices to support your business.\r\n    If you need a sample, please feel free to contact me,look forward to your response.\r\n</p>\r\n<img src=\"cid:content3.png\" alt=\"glove image\">\r\n<p>Best regards,<br>\r\nUna Huang\r\n</p>\r\n<p></p>\r\n<p>Sunwei Industrial Co., Ltd.<br>\r\nEmail: sales2@sunweiglove.com<br>\r\nTell & Whats app:+86-13751336933<br>\r\nECONOMIC DEVELOPMENT ZONE,CAOXIAN COUNTY, SHANDONG PROVINCE, P.R. CHINA.<br>\r\n    <span style=\"color: yellowgreen\">Please consider the environment before printing this e-mail.</span><br>\r\n</p>\r\n</div>\r\n</body>\r\n</html>\r\n', '[6]', 3, '2024-08-10 11:32:59', NULL);
INSERT INTO `email_content` VALUES (4, 'Latex Allergies? Our Gloves Offer a Safe Alternative', '<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Title</title>\r\n</head>\r\n<body style=\" font-family: Arial, sans-serif;\r\n            background-color: #f5f5f5;\r\n            color: #333333;\r\n            margin: 0;\r\n            padding: 20px;\">\r\n<div style=\"font-size: 16px;\r\n            line-height: 1.6;\r\n            background-color: #ffffff;\r\n            max-width: 600px;\r\n            padding: 20px;\r\n            border-radius: 8px;\r\n            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\">\r\n<p>Dear,<br>\r\nGreetings!\r\n</p>\r\n    <p>  I am Una Huang from <span style=\"font-weight: bold\">Sunwei Industrial Co.</span>, Ltd.We understand that latex allergies are a significant concern for many users. Our nitrile and vinyl gloves are latex-free, providing a safe alternative for those with latex sensitivities.\r\n</p>\r\n<p>we would love to provide you with samples to demonstrate our commitment to quality, as well as quotation for your evaluation. Looking forward to your contact.</p>\r\n<img src=\"cid:content4-1.png\" alt=\"glove image\">\r\n<img src=\"cid:content4-2.png\" alt=\"glove image\">\r\n<img src=\"cid:content4-3.png\" alt=\"glove image\">\r\n<p>Best regards,<br>\r\nUna Huang\r\n</p>\r\n<p></p>\r\n<p>Sunwei Industrial Co., Ltd.<br>\r\nEmail: sales2@sunweiglove.com<br>\r\nTell & Whats app:+86-13751336933<br>\r\nECONOMIC DEVELOPMENT ZONE,CAOXIAN COUNTY, SHANDONG PROVINCE, P.R. CHINA.<br>\r\n    <span style=\"color: yellowgreen\">Please consider the environment before printing this e-mail.</span><br>\r\n</p>\r\n</div>\r\n</body>\r\n</html>\r\n', '[7, 8, 9, 10]', 4, '2024-08-10 11:33:02', NULL);

-- ----------------------------
-- Table structure for email_task
-- ----------------------------
DROP TABLE IF EXISTS `email_task`;
CREATE TABLE `email_task`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮件地址',
  `content_id` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '邮件内容id',
  `is_replay` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否回复,0:未回复，1：已回复',
  `send_time` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '发送时间',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '邮件定时任务' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of email_task
-- ----------------------------

-- ----------------------------
-- Table structure for search_contact
-- ----------------------------
DROP TABLE IF EXISTS `search_contact`;
CREATE TABLE `search_contact`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮件地址',
  `phone` varchar(19) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '电话',
  `category` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类,1:手动,2:google',
  `keyword` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '关键词',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'url',
  `md5` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '验证urll唯一',
  `create_time` datetime NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '邮件定时任务' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of search_contact
-- ----------------------------
INSERT INTO `search_contact` VALUES (1, '', '', '2', 'disposable gloves contact email phone', 'https://shopping.medexpressgloves.com/', 'd41d8cd98f00b204e9800998ecf8427e', '2024-08-11 16:31:37', '2024-08-11 16:31:37');
INSERT INTO `search_contact` VALUES (2, '', '', '2', 'disposable gloves contact email phone', 'https://shopping.medexpressgloves.com/', 'aaa22613f5880c4f4c5768ffc2fb8c41', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (3, '', '', '2', 'disposable gloves contact email phone', 'https://www.mechanix.com/us-en/', 'f298c2b642cceeed2b9828966bd9dc5f', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (4, '', '', '2', 'disposable gloves contact email phone', 'https://ehs.berkeley.edu/glove-selection-guide', '8beb30f521e1c27a153ef276272fc84a', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (5, '', '', '2', 'disposable gloves contact email phone', 'https://www.libertysafety.com/contact/', 'ae2d0b0e25cff0b4a21113791730fbd6', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (6, '', '', '2', 'disposable gloves contact email phone', 'https://www.cypressmed.com/contact-us/', 'fe3c2b14510a944d78a0a459ff947b21', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (7, '', '', '2', 'disposable gloves contact email phone', 'https://pdihc.com/products/environment-of-care/super-sani-cloth-germicidal-disposable-wipe/', 'd864114e265cc9be4201f5f652ed0d67', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (8, '', '', '2', 'disposable gloves contact email phone', 'https://www.lifeguardgloves.com/', 'd481bdce911b1f0f610ad37b55dcb55f', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (9, '731847483@qq.com', '', '1', 'disposable gloves contact email phone', '', '731847483@qq.com', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (10, 'davezhou523@gmail.com', '', '1', 'disposable gloves contact email phone', '', 'davezhou523@gmail.com', '2024-08-11 16:32:02', NULL);
INSERT INTO `search_contact` VALUES (11, '271416962@qq.com', '', '1', 'disposable gloves contact email phone', '', '271416962@qq.com', '2024-08-11 16:32:02', NULL);

SET FOREIGN_KEY_CHECKS = 1;
