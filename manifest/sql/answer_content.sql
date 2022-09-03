/*
 Source Server Type    : MySQL
 Source Server Version : 50718
 File Encoding         : 65001
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for answer_content
-- ----------------------------
DROP TABLE IF EXISTS `answer_content`;
CREATE TABLE `answer_content`  (
  `answer_id` bigint(20) NOT NULL COMMENT '回答 id',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '回答内容',
  `meta` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '元数据',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`answer_id`) USING BTREE,
  CONSTRAINT `answer_content_ibfk_1` FOREIGN KEY (`answer_id`) REFERENCES `answer_index` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
