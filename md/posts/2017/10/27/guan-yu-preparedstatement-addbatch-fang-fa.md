---
title: "关于PreparedStatement.addBatch()方法"
categories: [ "java" ]
tags: [ "statement","preparedStatement" ]
draft: false
slug: "guan-yu-preparedstatement-addbatch-fang-fa"
date: "2017-10-27 13:59:00"
url: "/guan-yu-preparedstatement-addbatch-fang-fa.html"
---

<blockquote>Statement和PreparedStatement的区别就不多废话了,直接说PreparedStatement最重要的addbatch()结构的使用. 1.建立链接,(打电话拨号 )        Connection    connection =getConnection(); 2.不自动 Commit (瓜子不是一个一个吃,全部剥开放桌子上,然后一口舔了) con</blockquote>
<article>
<div id="article_content" class="article_content csdn-tracking-statistics" data-mod="popu_307" data-dsm="post">

Statement和PreparedStatement的区别就不多废话了,直接说PreparedStatement最重要的addbatch()结构的使用.

1.建立链接,(打电话拨号 )

Connection    connection =getConnection();

2.不自动 Commit (瓜子不是一个一个吃,全部剥开放桌子上,然后一口舔了)

connection.setAutoCommit(false);

3.预编译SQL语句,只编译一回哦,效率高啊.(发明一个剥瓜子的方法,以后不要总想怎么剥瓜子好.就这样剥.)
PreparedStatement statement = connection.prepareStatement("INSERT INTO TABLEX VALUES(?, ?)");

4.来一个剥一个,然后放桌子上

//记录1
statement.setInt(1, 1);
statement.setString(2, "Cujo");
statement.addBatch();

//记录2
statement.setInt(1, 2);
statement.setString(2, "Fred");
statement.addBatch();

//记录3
statement.setInt(1, 3);
statement.setString(2, "Mark");
statement.addBatch();

//批量执行上面3条语句. 一口吞了,很爽
int [] counts = statement.executeBatch();

//Commit it 咽下去,到肚子(DB)里面
connection.commit();

stmt.addBatch("update  TABLE1 set 题目="盛夏话足部保健1"   where id="3407"");
stmt.addBatch("update  TABLE1 set 题目="夏季预防中暑膳食1" where id="3408"");
stmt.addBatch("INSERT INTO  TABLE1  VALUES("11","12","13","","")");
stmt.addBatch("INSERT INTO  TABLE1  VALUES("12","12","13","","")");
stmt.addBatch("INSERT INTO  TABLE1  VALUES("13","12","13","","")");
stmt.addBatch("INSERT INTO  TABLE1  VALUES("14","12","13","","")");
stmt.addBatch("INSERT INTO  TABLE1  VALUES("15","12","13","","")");
stmt.addBatch("INSERT INTO  TABLE1  VALUES("16","12","13","","")");
stmt.addBatch("INSERT INTO  TABLE1  VALUES("17","12","13","","")");
stmt.addBatch("INSERT INTO  TABLE1  VALUES("18","12","13","","")");

int [] updateCounts=stmt.executeBatch();
cn.commit();

例如:

public static void execteBatch(Connection conn)throws Exception{
String sql1 = "delete from student where id =3 ";
String sql2 = "delete from student where id =5 ";
String sql3 = "delete from student where id =6 ";
String sql4 = "delete from student where id =7 ";
PreparedStatement pstmt = conn.prepareStatement(sql1);
pstmt.addBatch();
pstmt.addBatch(sql2);
pstmt.addBatch(sql3);
pstmt.addBatch(sql4);
pstmt.executeBatch();
};

</div>
</article>