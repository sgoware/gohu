/*
 Source Server Type    : MySQL
 Source Server Version : 50718
 File Encoding         : 65001
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment_index
-- ----------------------------
DROP TABLE IF EXISTS `comment_index`;
CREATE TABLE `comment_index`  (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `subject_id` bigint(20) NOT NULL COMMENT '评论主题 id',
  `user_id` bigint(20) NOT NULL COMMENT '评论者/回复者 id',
  `ip_loc` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论者/回复者 归属地',
  `root_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '评论 id (为 0 则是回复)',
  `comment_floor` int(11) NOT NULL DEFAULT 0 COMMENT '当前评论层数 (若是回复则为 0)',
  `comment_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '回复对象 id (若是评论则为 0)',
  `reply_floor` int(11) NOT NULL DEFAULT 0 COMMENT '当前回复层数 (若是评论则为 0)',
  `approve_count` int(11) NOT NULL DEFAULT 0 COMMENT '赞同数',
  `state` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态 (0-正常 1-隐藏)',
  `attrs` tinyint(4) NOT NULL COMMENT '属性 (待添加)',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `subject_id`(`subject_id`) USING BTREE,
  INDEX `root_id`(`root_id`) USING BTREE,
  CONSTRAINT `comment_index_ibfk_1` FOREIGN KEY (`subject_id`) REFERENCES `comment_subject` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
