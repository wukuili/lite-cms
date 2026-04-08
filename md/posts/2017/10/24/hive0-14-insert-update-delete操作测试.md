---
title: "hive0.14-insert、update、delete操作测试"
categories: [ "hive" ]
tags: [ "hive","insert","update","delete" ]
draft: false
slug: "hive0-14-insert-update-delete操作测试"
date: "2017-10-24 17:39:04"
url: "/hive0-14-insert-update-delete操作测试.html"
---

<table cellspacing="0" cellpadding="0">
<tbody>
<tr>
<td id="postmessage_76762" class="t_f"><strong>问题导读</strong><strong><span style="color: #ff0000;">1.测试insert报错，该如何解决？
2.hive delete和update报错，该如何解决？
3.什么情况下才允许delete和update？
</span></strong>

<img id="aimg_D99YS" class="zoom" src="https://lequ7.com/wp-content/uploads/2017/10/4.gif" alt="" width="500" height="35" border="0" />

首先用最普通的建表语句建一个表：
<div class="blockcode">
<div id="code_lJ1">
<ol>
 	<li>hive&gt;create table test(id int,name string)row format delimited fields terminated by ',';</li>
</ol>
</div>
<em>复制代码</em>

</div>
测试insert：
<div class="blockcode">
<div id="code_EAD">
<ol>
 	<li>insert into table test values (1,'row1'),(2,'row2');</li>
</ol>
</div>
<em>复制代码</em>

</div>
结果报错：
<div class="blockcode">
<div id="code_DvV">
<ol>
 	<li>java.io.FileNotFoundException: File does not exist: hdfs://127.0.0.1:9000/home/hadoop/git/hive/packaging/target/apache-hive-0.14.0-SNAPSHOT-bin/</li>
 	<li>apache-hive-0.14.0-SNAPSHOT-bin/lib/curator-client-2.6.0.jar</li>
 	<li>        at org.apache.hadoop.hdfs.DistributedFileSystem$17.doCall(DistributedFileSystem.java:1128)</li>
 	<li>        at org.apache.hadoop.hdfs.DistributedFileSystem$17.doCall(DistributedFileSystem.java:1120)</li>
 	<li>        at org.apache.hadoop.fs.FileSystemLinkResolver.resolve(FileSystemLinkResolver.java:81)</li>
 	<li>        at org.apache.hadoop.hdfs.DistributedFileSystem.getFileStatus(DistributedFileSystem.java:1120)</li>
 	<li>        at org.apache.hadoop.mapreduce.filecache.ClientDistributedCacheManager.getFileStatus(ClientDistributedCacheManager.java:288)</li>
 	<li>        at org.apache.hadoop.mapreduce.filecache.ClientDistributedCacheManager.getFileStatus(ClientDistributedCacheManager.java:224)</li>
 	<li>        at org.apache.hadoop.mapreduce.filecache.ClientDistributedCacheManager.determineTimestamps(ClientDistributedCacheManager.java:99)</li>
 	<li>        at org.apache.hadoop.mapreduce.filecache.ClientDistributedCacheManager.determineTimestampsAndCacheVisibilities(ClientDistributedCacheManager.java:57)</li>
 	<li>        at org.apache.hadoop.mapreduce.JobSubmitter.copyAndConfigureFiles(JobSubmitter.java:265)</li>
 	<li>        at org.apache.hadoop.mapreduce.JobSubmitter.copyAndConfigureFiles(JobSubmitter.java:301)</li>
 	<li>        at org.apache.hadoop.mapreduce.JobSubmitter.submitJobInternal(JobSubmitter.java:389)</li>
 	<li>        at org.apache.hadoop.mapreduce.Job$10.run(Job.java:1285)</li>
 	<li>        at org.apache.hadoop.mapreduce.Job$10.run(Job.java:1282)</li>
 	<li>        at java.security.AccessController.doPrivileged(Native Method)</li>
 	<li>        ......</li>
</ol>
</div>
<em>复制代码</em>

</div>
貌似往hdfs上找jar包了，小问题，直接把lib下的jar包上传到hdfs
<div class="blockcode">
<div id="code_Ih9">
<ol>
 	<li>hadoop fs -mkdir -p /home/hadoop/git/hive/packaging/target/apache-hive-0.14.0-SNAPSHOT-bin/apache-hive-0.14.0-SNAPSHOT-bin/lib/</li>
 	<li>hadoop fs -put $HIVE_HOME/lib/* /home/hadoop/git/hive/packaging/target/apache-hive-0.14.0-SNAPSHOT-bin/apache-hive-0.14.0-SNAPSHOT-bin/lib/</li>
 	<li></li>
</ol>
</div>
<em>复制代码</em>

</div>
接着运行insert，没有问题，接下来测试delete
<div class="blockcode">
<div id="code_iQB">
<ol>
 	<li>hive&gt;delete from test where id = 1;</li>
</ol>
</div>
<em>复制代码</em>

</div>
报错！：
FAILED: SemanticException [Error 10294]: Attempt to do update or delete using transaction manager that does not support these operations.
说是在使用的转换管理器不支持update跟delete操作。
原来要支持update操作跟delete操作，必须额外再配置一些东西，见：
<a href="https://cwiki.apache.org/confluence/display/Hive/Hive+Transactions#HiveTransactions-NewConfigurationParametersforTransactions" target="_blank" rel="noopener">https://cwiki.apache.org/conflue ... tersforTransactions</a>
根据提示配置hive-site.xml:
<div class="blockcode">
<div id="code_uhK">
<ol>
 	<li>    hive.support.concurrency – true</li>
 	<li>    hive.enforce.bucketing – true</li>
 	<li>    hive.exec.dynamic.partition.mode – nonstrict</li>
 	<li>    hive.txn.manager – org.apache.hadoop.hive.ql.lockmgr.DbTxnManager</li>
 	<li>    hive.compactor.initiator.on – true</li>
 	<li>    hive.compactor.worker.threads – 1</li>
</ol>
</div>
<em>复制代码</em>

</div>
配置完以为能够顺利运行了，谁知开始报下面这个错误：
<div class="blockcode">
<div id="code_lR6">
<ol>
 	<li>FAILED: LockException [Error 10280]: Error communicating with the metastore</li>
</ol>
</div>
<em>复制代码</em>

</div>
与元数据库出现了问题，修改log为DEBUG查看具体错误：
<div class="blockcode">
<div id="code_t9b">
<ol>
 	<li>4-11-04 14:20:14,367 DEBUG [Thread-8]: txn.CompactionTxnHandler (CompactionTxnHandler.java:findReadyToClean(265)) - Going to execute query &lt;select cq_id,</li>
 	<li>cq_database, cq_table, cq_partition, cq_type, cq_run_as from COMPACTION_QUEUE where cq_state = 'r'&gt;</li>
 	<li>2014-11-04 14:20:14,367 ERROR [Thread-8]: txn.CompactionTxnHandler (CompactionTxnHandler.java:findReadyToClean(285)) - Unable to select next element for cleaning,</li>
 	<li>Table 'hive.COMPACTION_QUEUE' doesn't exist</li>
 	<li>2014-11-04 14:20:14,367 DEBUG [Thread-8]: txn.CompactionTxnHandler (CompactionTxnHandler.java:findReadyToClean(287)) - Going to rollback</li>
 	<li>2014-11-04 14:20:14,368 ERROR [Thread-8]: compactor.Cleaner (Cleaner.java:run(143)) - Caught an exception in the main loop of compactor cleaner, MetaException(message</li>
 	<li>:Unable to connect to transaction database com.mysql.jdbc.exceptions.jdbc4.MySQLSyntaxErrorException: Table 'hive.COMPACTION_QUEUE' doesn't exist</li>
 	<li>    at sun.reflect.GeneratedConstructorAccessor19.newInstance(Unknown Source)</li>
 	<li>    at sun.reflect.DelegatingConstructorAccessorImpl.newInstance(DelegatingConstructorAccessorImpl.java:45)</li>
 	<li>    at java.lang.reflect.Constructor.newInstance(Constructor.java:526)</li>
 	<li>    at com.mysql.jdbc.Util.handleNewInstance(Util.java:409)</li>
</ol>
</div>
<em>复制代码</em>

</div>
在元数据库中找不到COMPACTION_QUEUE这个表，赶紧去mysql中查看，确实没有这个表。怎么会没有这个表呢？找了很久都没找到什么原因，查源码吧。
在org.apache.hadoop.hive.metastore.txn下的TxnDbUtil类中找到了建表语句，顺藤摸瓜，找到了下面这个方法会调用建表语句：
<div class="blockcode">
<div id="code_OzP">
<ol>
 	<li>private void checkQFileTestHack() {</li>
 	<li>    boolean hackOn = HiveConf.getBoolVar(conf, HiveConf.ConfVars.HIVE_IN_TEST) ||</li>
 	<li>        HiveConf.getBoolVar(conf, HiveConf.ConfVars.HIVE_IN_TEZ_TEST);</li>
 	<li>    if (hackOn) {</li>
 	<li>      LOG.info("Hacking in canned values for transaction manager");</li>
 	<li>      // Set up the transaction/locking db in the derby metastore</li>
 	<li>      TxnDbUtil.setConfValues(conf);</li>
 	<li>      try {</li>
 	<li>        TxnDbUtil.prepDb();</li>
 	<li>      } catch (Exception e) {</li>
 	<li>        // We may have already created the tables and thus don't need to redo it.</li>
 	<li>        if (!e.getMessage().contains("already exists")) {</li>
 	<li>          throw new RuntimeException("Unable to set up transaction database for" +</li>
 	<li>              " testing: " + e.getMessage());</li>
 	<li>        }</li>
 	<li>      }</li>
 	<li>    }</li>
 	<li>  }</li>
</ol>
</div>
<em>复制代码</em>

</div>
什么意思呢，就是说要运行建表语句还有一个条件：HIVE_IN_TEST或者HIVE_IN_TEZ_TEST.只有在测试环境中才能用delete，update操作，也可以理解，毕竟还没有开发完全。
终于找到原因，解决方法也很简单：在hive-site.xml中添加下面的配置：
<div class="blockcode">
<div id="code_yyZ">
<ol>
 	<li>&lt;property&gt;</li>
 	<li>&lt;name&gt;hive.in.test&lt;/name&gt;</li>
 	<li>&lt;value&gt;true&lt;/value&gt;</li>
 	<li>&lt;/property&gt;</li>
</ol>
</div>
<em>复制代码</em>

</div>
OK,再重新启动服务，再运行delete：
<div class="blockcode">
<div id="code_AjG">
<ol>
 	<li>hive&gt;delete from test where id = 1;</li>
</ol>
</div>
<em>复制代码</em>

</div>
又报错：
<div class="blockcode">
<div id="code_PNH">
<ol>
 	<li>FAILED: SemanticException [Error 10297]: Attempt to do update or delete on table default.test that does not use an AcidOutputFormat or is not bucketed</li>
</ol>
</div>
<em>复制代码</em>

</div>
说是要进行delete操作的表test不是AcidOutputFormat或没有分桶。估计是要求输出是AcidOutputFormat然后必须分桶
网上查到确实如此，而且目前只有ORCFileformat支持AcidOutputFormat，不仅如此建表时必须指定参数('transactional' = true)。感觉太麻烦了。。。。
于是按照网上示例建表：
<div class="blockcode">
<div id="code_X5i">
<ol>
 	<li>hive&gt;create table test(id int ,name string )clustered by (id) into 2 buckets stored as orc TBLPROPERTIES('transactional'='true');</li>
</ol>
</div>
<em>复制代码</em>

</div>
insert
<div class="blockcode">
<div id="code_kxQ">
<ol>
 	<li>hive&gt;insert into table test values (1,'row1'),(2,'row2'),(3,'row3');</li>
</ol>
</div>
<em>复制代码</em>

</div>
delete
<div class="blockcode">
<div id="code_p5K">
<ol>
 	<li>hive&gt;delete from test where id = 1;</li>
</ol>
</div>
<em>复制代码</em>

</div>
update
<div class="blockcode">
<div id="code_qO7">
<ol>
 	<li>hive&gt;update test set name = 'Raj' where id = 2;</li>
</ol>
</div>
<em>复制代码</em>

</div>
OK!全部顺利运行，不过貌似效率太低了，基本都要30s左右，估计应该可以优化，再研究研究</td>
</tr>
</tbody>
</table>
<!--more-->