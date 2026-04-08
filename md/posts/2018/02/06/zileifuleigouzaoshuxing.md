---
title: "子类使用父类构造方法设置自己的属性"
categories: [ "java" ]
tags: [ "继承","java" ]
draft: false
slug: "zileifuleigouzaoshuxing"
date: "2018-02-06 15:36:15"
url: "/zileifuleigouzaoshuxing.html"
---

首先需要知道执行流程，当new一个子类的实例时，执行顺序为 父类构造方法---》子类属性赋值

所以可以通过对子类属性设置set和get方法来实现对子类属性的赋值，请看以下代码

父类：
<pre> [code lang="java"]
 import com.taiji.tr.errorlog.MyLogger;

import java.util.Set;

/**
 * @Author: liyj
 * @Description:
* @Date:Created in 2018/1/25
 * @Modified by :
 */
 public abstract class ProcedureSql {

CompProperDataBean compProperDataBean;

public ProcedureSql(CompProperDataBean compProperDataBean) {
 this.compProperDataBean = compProperDataBean;
 init();
 }

public CompProperDataBean getCompProperDataBean() {
 return compProperDataBean;
 }

public abstract void init();
 }
 [/code]</pre>
子类：

[code lang="java"]
package com.taiji.tr.bean;

import com.taiji.tr.errorlog.MyLogger;

/**
* @Author: liyj
* @Description: &amp;lt;p/&amp;gt;
* @Date:Created in 2018/1/31
* @Modified by :
*/
public class HiveProcedureSql extends ProcedureSql {
public String tempSqlStr;

public String getTempSqlStr() {
return tempSqlStr;
}

public void setTempSqlStr(String tempSqlStr) {
this.tempSqlStr = tempSqlStr;
}

public HiveProcedureSql(CompProperDataBean compProperDataBean) {
super(compProperDataBean);
}

@Override
public void init() {
//判断是否需要创建分区
if (compProperDataBean.isPartition()) {//创建分区
setTempSqlStr(&quot;' INTO table &quot; + compProperDataBean.getTable_name() + &quot; PARTITION (&quot;);
} else {
setTempSqlStr(&quot;' into table &quot; + compProperDataBean.getTable_name());
}
MyLogger.info(&quot;sql语句&quot;, compProperDataBean.getIdentify(), tempSqlStr);
}
}
[/code]