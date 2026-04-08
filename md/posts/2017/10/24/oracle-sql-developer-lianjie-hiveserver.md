---
title: "Oracle SQL Developer连接HiveServer"
categories: [ "hive" ]
tags: [ "hive","sql Developer" ]
draft: false
slug: "oracle-sql-developer-lianjie-hiveserver"
date: "2017-10-24 15:25:00"
url: "/oracle-sql-developer-lianjie-hiveserver.html"
---

<blockquote>Oracle SQL Developer从http://www.oracle.com/technetwork/developer-tools/sql-developer/downloads/index.html下载SQL Developer 4.1.5，并解压；从http://www.cloudera.com/downloads/connectors/hive/jdbc/2-5-15.html下载Hive JDBC Driver for Oracle SQL Developer，并解压，进入解压后的目录，将Cloudera_HiveJDBC4_2.5.15.1040.zip解压。打开sqldeveloper.exe，点击”工具”–&gt;“首选项”,在”数据库”–&gt;”第三方JDBC驱动”中，添加Hive JDBC驱动：添加后重启sqldeveloper。再次打开sqldeveloper后，点击”新建连接”之后，多了”Hive”数据库：连接Hive:</blockquote>
<h2>Oracle SQL Developer</h2>
从http://www.oracle.com/technetwork/developer-tools/sql-developer/downloads/index.html下载SQL Developer 4.1.5，并解压；

从http://www.cloudera.com/downloads/connectors/hive/jdbc/2-5-15.html下载Hive JDBC Driver for Oracle SQL Developer，并解压，进入解压后的目录，将Cloudera_HiveJDBC4_2.5.15.1040.zip解压。

打开sqldeveloper.exe，点击”工具”–&gt;“首选项”,在”数据库”–&gt;”第三方JDBC驱动”中，添加Hive JDBC驱动：

<img class="aligncenter" src="https://lequ7.com/wp-content/uploads/2017/10/20160928-4.jpg" alt="hive" width="1330" height="748" />

添加后重启sqldeveloper。

再次打开sqldeveloper后，点击”新建连接”之后，多了”Hive”数据库：

<img class="aligncenter" src="https://lequ7.com/wp-content/uploads/2017/10/20160928-6.jpg" alt="hive" width="741" height="448" />

连接Hive:

<img class="alignnone size-full" src="https://lequ7.com/wp-content/uploads/2017/10/20160928-5.jpg" alt="" />