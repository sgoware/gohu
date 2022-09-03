/*
 Source Server Type    : MySQL
 Source Server Version : 50718
 File Encoding         : 65001
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment_subject
-- ----------------------------
DROP TABLE IF EXISTS `comment_subject`;
CREATE TABLE `comment_subject`  (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `obj_type` tinyint(4) NOT NULL COMMENT '对象类型 (0-默认对象 1-回答 2-文章)',
  `obj_id` bigint(20) NOT NULL COMMENT '对象 id',
  `count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论+回复总数',
  `root_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论总数',
  `state` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态 (0-正常 1-隐藏)',
  `attrs` tinyint(4) UNSIGNED NULL DEFAULT 0 COMMENT '属性 (待添加)',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `obj_id`(`obj_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
