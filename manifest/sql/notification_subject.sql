/*
 Source Server Type    : MySQL
 Source Server Version : 50718
 File Encoding         : 65001
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for notification_subject
-- ----------------------------
DROP TABLE IF EXISTS `notification_subject`;
CREATE TABLE `notification_subject`  (
  `id` bigint(20) NOT NULL COMMENT '主键',
  `user_id` bigint(20) NOT NULL COMMENT '用户 ID',
  `message_type` tinyint(4) NOT NULL DEFAULT 0 COMMENT '对象类型 (0-全部 1-关注我的 2-赞同与喜欢 3-评论与回复 4-关注的人 5-问题回答)',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `obj_type`(`message_type`) USING BTREE,
  CONSTRAINT `notification_subject_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user_subject` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
