DROP TABLE IF EXISTS `balance_records`;
CREATE TABLE `balance_records` (
    `id` bigint(20) unsigned NOT NULL  AUTO_INCREMENT COMMENT 'id',
    `user_id` bigint(20)  NOT NULL DEFAULT 0 COMMENT 'user_id',
    `balance` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT 'current balance',
    `deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'the record is deleted or not',
    `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'record created time',
    `updated_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'record updated time',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='balance record';
