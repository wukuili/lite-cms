---
title: "log4j: log4j:ERROR Attempted to append to closed appender named [stdout] "
categories: [ "java" ]
tags: [ "log4j","appender named" ]
draft: false
slug: "log4j-log4jerror-attempted-to-append-to-closed-appender-named-stdout"
date: "2017-10-26 09:55:46"
url: "/log4j-log4jerror-attempted-to-append-to-closed-appender-named-stdout.html"
---

<blockquote>log4j: log4j:ERROR Attempted to append to closed appender named [stdout]</blockquote>
发给现场用的jar包却发现报上面的错，可是在家里电脑上却无法复现出来，果断搜索，才发现原来是多个log4j.properties文件引起的。删除其中一个就行了。

另外如果你那里也遇到了相同的问题，但是只有一个log4j.properties,那么你就要看看自己的程序里面是不是有两种方式获取log4j对象（一种是直接放到了classpath文件夹，另一种是通过PropertyConfigurator.configure(..)这种方式）