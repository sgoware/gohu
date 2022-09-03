/*
 Source Server Type    : MySQL
 Source Server Version : 50718
 File Encoding         : 65001
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for question_subject
-- ----------------------------
DROP TABLE IF EXISTS `question_subject`;
CREATE TABLE `question_subject`  (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `user_id` bigint(20) NOT NULL COMMENT '提问者 id',
  `ip_loc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '提问者 IP 归属地',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '问题标题',
  `topic` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '问题主题',
  `tag` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '问题标签',
  `sub_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '关注总数',
  `answer_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '回答总数',
  `view_count` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览总数',
  `state` tinyint(4) UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态 (0-正常 1-隐藏)',
  `attrs` tinyint(4) NULL DEFAULT NULL COMMENT '属性 (待添加)',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
