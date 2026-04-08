---
title: "python练习时出现module &#039;re&#039; has no attribute &#039;match&#039;错误"
categories: [ "日常" ]
tags: [ "python","re","module" ]
draft: false
slug: "pythonlianxishichuxianmodule-re-has-no-attribute-matchcuowu"
date: "2018-04-06 17:06:21"
url: "/pythonlianxishichuxianmodule-re-has-no-attribute-matchcuowu.html"
---

练习python中正则表达式，使用import re，并且写了最简单的匹配，运行py文件，却提示module 're' has no attribute 'match'错误。

仔细查看原因，原来是自己把python package写成了re了。

导致python在执行import re时，以为我们要他导入我自己定义的package呢，将包名字改成retest，重新运行py文件，成功运行。