---
title: "kafka [B cannot be cast to java.lang.String"
categories: [ "Kafka" ]
tags: [ "Kafka","StringEncoder","DefaultEncoder" ]
draft: false
slug: "kafka-b-cannot-be-cast-to-java-lang-string"
date: "2017-10-24 14:54:21"
url: "/kafka-b-cannot-be-cast-to-java-lang-string.html"
---

自己写了一个Kafka发送消息的demo，但是发送消息的时候，却报了kafka [B cannot be cast to java.lang.String的错误，后来找到了解决办法

原来是因为在定义config文件的时候，针对
<pre class="lang-scala prettyprint prettyprinted"><code><span class="str">serializer.class</span></code>
部分，错误的当成了StringEncoder，其实修改成默认的encoder就行了</pre>
<pre class="lang-scala prettyprint prettyprinted"><code><span class="pln">props</span><span class="pun">.</span><span class="pln">put</span><span class="pun">(</span><span class="str">"serializer.class"</span><span class="pun">,</span> <span class="str">"kafka.serializer.DefaultEncoder"</span><span class="pun">);</span></code></pre>