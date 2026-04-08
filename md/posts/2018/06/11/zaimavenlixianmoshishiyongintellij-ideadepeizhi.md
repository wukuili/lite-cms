---
title: "在maven离线模式使用intellij idea的配置"
categories: [ "java","日常" ]
tags: [ "idea","intellij idea","maven","offline","离线模式" ]
draft: false
slug: "zaimavenlixianmoshishiyongintellij-ideadepeizhi"
date: "2018-06-11 19:44:48"
url: "/zaimavenlixianmoshishiyongintellij-ideadepeizhi.html"
---

今天在内网使用intellij idea，把maven仓库也拷贝了过来，但是使用maven install 命令的时候，却还是从私服里面下载。（内网无法连接到任何私服）导致install的时候报缺少依赖的错误。

经过各种搜索，终于发现了解决办法。

<!--more-->

需要修改下maven的setting.xml文件。

找打自己使用的setting.xml文件，找到仓库私服连接的配置的行数，将私服的地址换成file：地址。如下


[code lang="xml"]

&lt;repository&gt;
&lt;id&gt;central&lt;/id&gt;
&lt;url&gt;file://D:\mavenrepo&lt;/url&gt;
&lt;/repository&gt;

[/code]


之后重新导入依赖，然后可以成功执行maven install 。