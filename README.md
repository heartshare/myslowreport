# myslowreport

### 需要
1. mysql服务器上配置rsync用于同步慢日志,同时部署切割慢日志脚本cut_myslow.sh
2. 部署一台mysql数据库用于导入慢日志
3. 运行报告程序的机器部署perconatoolkit工具包

### 执行流程
1. rsync同步mysql慢日志到报告程序的机器
2. 使用pt-query-digest导入mysql慢日志到mysql数据库
3. 查询导入的慢日志记录后生成html形式报告发送至相关人员

### 注意
pt-query-digest 导入时会缺少 db_max user_max 两个字段, 因此需要事先建好表, 表结构如下:

表名根据实际情况修改,这个有对应到conf下的配置的

```Java
CREATE TABLE `myslow_history_10_120_3306` (
  `db_max` varchar(100) DEFAULT NULL,
  `user_max` varchar(100) DEFAULT NULL,
  `checksum` bigint(20) unsigned NOT NULL,
  `sample` text NOT NULL,
  `ts_min` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `ts_max` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `ts_cnt` float DEFAULT NULL,
  `Query_time_sum` float DEFAULT NULL,
  `Query_time_min` float DEFAULT NULL,
  `Query_time_max` float DEFAULT NULL,
  `Query_time_pct_95` float DEFAULT NULL,
  `Query_time_stddev` float DEFAULT NULL,
  `Query_time_median` float DEFAULT NULL,
  `Lock_time_sum` float DEFAULT NULL,
  `Lock_time_min` float DEFAULT NULL,
  `Lock_time_max` float DEFAULT NULL,
  `Lock_time_pct_95` float DEFAULT NULL,
  `Lock_time_stddev` float DEFAULT NULL,
  `Lock_time_median` float DEFAULT NULL,
  `Rows_sent_sum` float DEFAULT NULL,
  `Rows_sent_min` float DEFAULT NULL,
  `Rows_sent_max` float DEFAULT NULL,
  `Rows_sent_pct_95` float DEFAULT NULL,
  `Rows_sent_stddev` float DEFAULT NULL,
  `Rows_sent_median` float DEFAULT NULL,
  `Rows_examined_sum` float DEFAULT NULL,
  `Rows_examined_min` float DEFAULT NULL,
  `Rows_examined_max` float DEFAULT NULL,
  `Rows_examined_pct_95` float DEFAULT NULL,
  `Rows_examined_stddev` float DEFAULT NULL,
  `Rows_examined_median` float DEFAULT NULL,
  `Rows_affected_sum` float DEFAULT NULL,
  `Rows_affected_min` float DEFAULT NULL,
  `Rows_affected_max` float DEFAULT NULL,
  `Rows_affected_pct_95` float DEFAULT NULL,
  `Rows_affected_stddev` float DEFAULT NULL,
  `Rows_affected_median` float DEFAULT NULL,
  `Rows_read_sum` float DEFAULT NULL,
  `Rows_read_min` float DEFAULT NULL,
  `Rows_read_max` float DEFAULT NULL,
  `Rows_read_pct_95` float DEFAULT NULL,
  `Rows_read_stddev` float DEFAULT NULL,
  `Rows_read_median` float DEFAULT NULL,
  `Merge_passes_sum` float DEFAULT NULL,
  `Merge_passes_min` float DEFAULT NULL,
  `Merge_passes_max` float DEFAULT NULL,
  `Merge_passes_pct_95` float DEFAULT NULL,
  `Merge_passes_stddev` float DEFAULT NULL,
  `Merge_passes_median` float DEFAULT NULL,
  `InnoDB_IO_r_ops_min` float DEFAULT NULL,
  `InnoDB_IO_r_ops_max` float DEFAULT NULL,
  `InnoDB_IO_r_ops_pct_95` float DEFAULT NULL,
  `InnoDB_IO_r_ops_stddev` float DEFAULT NULL,
  `InnoDB_IO_r_ops_median` float DEFAULT NULL,
  `InnoDB_IO_r_bytes_min` float DEFAULT NULL,
  `InnoDB_IO_r_bytes_max` float DEFAULT NULL,
  `InnoDB_IO_r_bytes_pct_95` float DEFAULT NULL,
  `InnoDB_IO_r_bytes_stddev` float DEFAULT NULL,
  `InnoDB_IO_r_bytes_median` float DEFAULT NULL,
  `InnoDB_IO_r_wait_min` float DEFAULT NULL,
  `InnoDB_IO_r_wait_max` float DEFAULT NULL,
  `InnoDB_IO_r_wait_pct_95` float DEFAULT NULL,
  `InnoDB_IO_r_wait_stddev` float DEFAULT NULL,
  `InnoDB_IO_r_wait_median` float DEFAULT NULL,
  `InnoDB_rec_lock_wait_min` float DEFAULT NULL,
  `InnoDB_rec_lock_wait_max` float DEFAULT NULL,
  `InnoDB_rec_lock_wait_pct_95` float DEFAULT NULL,
  `InnoDB_rec_lock_wait_stddev` float DEFAULT NULL,
  `InnoDB_rec_lock_wait_median` float DEFAULT NULL,
  `InnoDB_queue_wait_min` float DEFAULT NULL,
  `InnoDB_queue_wait_max` float DEFAULT NULL,
  `InnoDB_queue_wait_pct_95` float DEFAULT NULL,
  `InnoDB_queue_wait_stddev` float DEFAULT NULL,
  `InnoDB_queue_wait_median` float DEFAULT NULL,
  `InnoDB_pages_distinct_min` float DEFAULT NULL,
  `InnoDB_pages_distinct_max` float DEFAULT NULL,
  `InnoDB_pages_distinct_pct_95` float DEFAULT NULL,
  `InnoDB_pages_distinct_stddev` float DEFAULT NULL,
  `InnoDB_pages_distinct_median` float DEFAULT NULL,
  `QC_Hit_cnt` float DEFAULT NULL,
  `QC_Hit_sum` float DEFAULT NULL,
  `Full_scan_cnt` float DEFAULT NULL,
  `Full_scan_sum` float DEFAULT NULL,
  `Full_join_cnt` float DEFAULT NULL,
  `Full_join_sum` float DEFAULT NULL,
  `Tmp_table_cnt` float DEFAULT NULL,
  `Tmp_table_sum` float DEFAULT NULL,
  `Tmp_table_on_disk_cnt` float DEFAULT NULL,
  `Tmp_table_on_disk_sum` float DEFAULT NULL,
  `Filesort_cnt` float DEFAULT NULL,
  `Filesort_sum` float DEFAULT NULL,
  `Filesort_on_disk_cnt` float DEFAULT NULL,
  `Filesort_on_disk_sum` float DEFAULT NULL,
  PRIMARY KEY (`checksum`,`ts_min`,`ts_max`),
  KEY `idx_checksum` (`checksum`) USING BTREE,
  KEY `idx_query_time_max` (`Query_time_max`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

哎, 浮点数是个罪孽, 修改部分字段类型为decimal 或者 bigint
```Java
ALTER TABLE myslow_history_10_120_3306
MODIFY COLUMN ts_cnt bigint(20) unsigned DEFAULT '0',

MODIFY COLUMN Query_time_sum decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Query_time_min decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Query_time_max decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Query_time_pct_95 decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Query_time_stddev decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Query_time_median decimal(12,9) DEFAULT '0.0',

MODIFY COLUMN Lock_time_sum decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Lock_time_min decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Lock_time_max decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Lock_time_pct_95 decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Lock_time_stddev decimal(12,9) DEFAULT '0.0',
MODIFY COLUMN Lock_time_median decimal(12,9) DEFAULT '0.0',

MODIFY COLUMN Rows_sent_sum bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_sent_min bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_sent_max bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_sent_pct_95 bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_sent_stddev bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_sent_median bigint(20) unsigned DEFAULT '0',

MODIFY COLUMN Rows_examined_sum bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_examined_min bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_examined_max bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_examined_pct_95 bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_examined_stddev bigint(20) unsigned DEFAULT '0',
MODIFY COLUMN Rows_examined_median bigint(20) unsigned DEFAULT '0'
;
```

增加增长率统计表
```Java
CREATE TABLE `myslow_history_grow_rate` (
	`Id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Id(自增量)',
	`MyInsName` VARCHAR(64) NOT NULL COMMENT 'MySQL实例名称',
	`StatDate` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '统计日期',
	`YesterdayTotal` int(11) NOT NULL DEFAULT '0' COMMENT '昨天的慢日志总数',
	`BeforeYesterdayTotal` int(11) NOT NULL DEFAULT '0' COMMENT '前天的慢日志总数',
	`BasisTotal` int(11) NOT NULL DEFAULT '0' COMMENT '上周此日的慢日志总数',
	`YesterdayUniq` int(11) NOT NULL DEFAULT '0' COMMENT '昨天的慢日志独立数',
	`BeforeYesterdayUniq` int(11) NOT NULL DEFAULT '0' COMMENT '前天的慢日志独立数',
	`BasisUniq` int(11) NOT NULL DEFAULT '0' COMMENT '上周此日的慢日志独立数',
	`TotalChainRate` DECIMAL(10,2) COMMENT '总慢日志语句数环比增长',
	`TotalBasisRate` DECIMAL(10,2) COMMENT '总慢日志语句数同比增长',
	`UniqChainRate` DECIMAL(10,2) COMMENT '独立慢日志语句数环比增长',
	`UniqBasisRate` DECIMAL(10,2) COMMENT '独立慢日志语句数同比增长',
	PRIMARY KEY (`Id`),
	UNIQUE KEY `idx_MyInsName_StatDate` (`MyInsName`,`StatDate`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```
