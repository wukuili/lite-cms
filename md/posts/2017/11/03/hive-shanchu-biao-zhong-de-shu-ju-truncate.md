---
title: "hive--删除表中的数据truncate"
categories: [ "hive" ]
tags: [ "hive","truncate" ]
draft: false
slug: "hive-shanchu-biao-zhong-de-shu-ju-truncate"
date: "2017-11-03 13:56:00"
url: "/hive-shanchu-biao-zhong-de-shu-ju-truncate.html"
---

delect:用于删除特定行条件,你可以从给定表中删除所有的行

TRUNCATE:truncate用于删除所有的行，这个行为在hive元存储删除数据是不可逆的

DROP:删除hive中的表

truncate 不能删除外部表！因为外部表里的数据并不是存放在Hive Meta store中

truncate:

truncate table table_name;

例子：

truncate table employees;