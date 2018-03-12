# myslowreport

### 需要
1. mysql服务器上配置rsync用于同步慢日志
2. 部署一台mysql数据库用于导入慢日志
3. 运行报告程序的机器部署perconatoolkit工具包

### 执行流程
1. rsync同步mysql慢日志到报告程序的机器
2. 使用pt-query-digest导入mysql慢日志到mysql数据库
3. 查询导入的慢日志记录后生成html形式报告发送至相关人员
