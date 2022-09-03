/*
 Source Server Type    : MySQL
 Source Server Version : 50718
 File Encoding         : 65001
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for answer_index
-- ----------------------------
DROP TABLE IF EXISTS `answer_index`;
CREATE TABLE `answer_index`  (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `question_id` bigint(20) NOT NULL COMMENT '问题 id',
  `user_id` bigint(20) NOT NULL COMMENT '回答者 id',
  `ip_loc` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '回答者 IP 归属地',
  `approve_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '赞同数',
  `like_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '喜欢数',
  `collect_count` int(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '收藏数',
  `state` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态 (0-正常 1-隐藏)',
  `attrs` tinyint(4) UNSIGNED NULL DEFAULT 0 COMMENT '属性 (待添加)',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `question_id`(`question_id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  CONSTRAINT `answer_index_ibfk_1` FOREIGN KEY (`question_id`) REFERENCES `question_subject` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
