/*
 Source Server Type    : MySQL
 Source Server Version : 50718
 File Encoding         : 65001
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_collection
-- ----------------------------
DROP TABLE IF EXISTS `user_collection`;
CREATE TABLE `user_collection`  (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `user_id` bigint(20) NOT NULL COMMENT '用户 id',
  `collect_type` tinyint(4) NOT NULL COMMENT '收藏类型 (1-喜欢 2-赞同 3-收藏 4-关注)',
  `obj_type` tinyint(4) NOT NULL COMMENT '对象类型 喜欢-(1-回答 2-文章) 赞同-(1-回答 2-文章 3-文章) 收藏-(0-默认对象 1-回答 2-文章) 关注-(1-用户 2-话题 3-专栏 4-问题 5-收藏夹)',
  `obj_id` bigint(20) NOT NULL COMMENT '对象 id',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `obj_type`(`obj_type`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `collect_type`(`collect_type`) USING BTREE,
  CONSTRAINT `user_collection_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user_subject` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
