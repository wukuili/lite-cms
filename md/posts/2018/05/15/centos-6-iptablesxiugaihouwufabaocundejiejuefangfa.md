---
title: "centos 6 iptables修改后无法保存的解决方法"
categories: [ "linux" ]
tags: [ "centos","ssh","iptables" ]
draft: false
slug: "centos-6-iptablesxiugaihouwufabaocundejiejuefangfa"
date: "2018-05-15 09:11:23"
url: "/centos-6-iptablesxiugaihouwufabaocundejiejuefangfa.html"
---

自己弄了一个vultur的便宜vps，主要用来爬梯子，但是发现频繁的遭受暴力密码破解的攻击（里面啥都没有，破解个球。。。。）

于是改端口，禁root登录，生成rsa文件，添加防火墙规则，重启ssh，可是发现仍然登不上去。

竟然发现知识用了防火墙的reload，没有用restart命令。。。。。

所以用vi改完 /etc/sysconfig/iptables 后，保存退出，然后启用

service iptables restart

重新登录，ok。

然后世界清静了。