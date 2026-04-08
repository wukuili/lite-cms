---
title: "教程：如何在vultur centos 7中部署Google BBR"
categories: [ "linux" ]
tags: [ "centos","centos7","vultur","bbr" ]
draft: false
slug: "jiaochengruhezaivultur-centos-7zhongbushugoogle-bbr"
date: "2018-03-30 08:28:35"
url: "/jiaochengruhezaivultur-centos-7zhongbushugoogle-bbr.html"
---

BBR（Bottleneck Bandwidth and RTT）是一种新的拥塞控制算法，由谷歌贡献给Linux内核TCP栈。使用BBR后，Linux服务器可以显着提高吞吐量并减少连接延迟。此外，部署BBR很容易，因为此算法仅需要发送方的更新，不需要网络或接收方。

在本文中，我将向您展示如何在Vultr CentOS 7 KVM服务器实例上部署BBR。

<!--more-->
<h3>步骤1：使用ELRepo RPM存储库升级内核</h3>
为了使用BBR，您需要将CentOS 7机器的内核升级到最新（本文以4.9.0版本为例）。您可以使用ELRepo RPM存储库轻松完成此操作。

在升级之前，你可以看看当前的内核：
<pre><code>uname -r
</code></pre>
该命令应输出一个类似于以下内容的字符串：
<pre><code>3.10.0-514.2.2.el7.x86_64
</code></pre>
如您所见，当前内核为3.10.0。

安装ELRepo回购：
<pre><code>sudo rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
sudo rpm -Uvh http://www.elrepo.org/elrepo-release-7.0-2.el7.elrepo.noarch.rpm
</code></pre>
使用ELRepo回购安装4.9.0内核：
<pre><code>sudo yum --enablerepo=elrepo-kernel install kernel-ml -y
</code></pre>
确认结果：
<pre><code>rpm -qa | grep kernel
</code></pre>
如果安装成功，您应该<code>kernel-ml-4.9.0-1.el7.elrepo.x86_64</code>在输出列表中看到：
<pre><code>kernel-ml-4.9.0-1.el7.elrepo.x86_64
kernel-3.10.0-514.el7.x86_64
kernel-tools-libs-3.10.0-514.2.2.el7.x86_64
kernel-tools-3.10.0-514.2.2.el7.x86_64
kernel-3.10.0-514.2.2.el7.x86_64
</code></pre>
现在，您需要通过设置默认的grub2启动项来启用4.9.0内核。

显示grub2菜单中的所有条目：
<pre><code>sudo egrep ^menuentry /etc/grub2.cfg | cut -f 2 -d \'
</code></pre>
结果应该类似于：
<pre><code>CentOS Linux 7 Rescue a0cbf86a6ef1416a8812657bb4f2b860 (4.9.0-1.el7.elrepo.x86_64)
CentOS Linux (4.9.0-1.el7.elrepo.x86_64) 7 (Core)
CentOS Linux (3.10.0-514.2.2.el7.x86_64) 7 (Core)
CentOS Linux (3.10.0-514.el7.x86_64) 7 (Core)
CentOS Linux (0-rescue-bf94f46c6bd04792a6a42c91bae645f7) 7 (Core)
</code></pre>
由于行数从<code>0</code>第二行开始，并且4.9.0内核条目位于第二行，因此将缺省引导条目设置为<code>1</code>：
<pre><code>sudo grub2-set-default 1
</code></pre>
重新启动系统：
<pre><code>sudo shutdown -r now
</code></pre>
当服务器恢复联机时，请重新登录并重新运行uname命令以确认您使用的是正确的内核：
<pre><code>uname -r
</code></pre>
你应该看到如下结果：
<pre><code>4.9.0-1.el7.elrepo.x86_64
</code></pre>
<h3>步骤2：启用BBR</h3>
为了启用BBR算法，您需要修改<code>sysctl</code>配置，如下所示：
<pre><code>echo 'net.core.default_qdisc=fq' | sudo tee -a /etc/sysctl.conf
echo 'net.ipv4.tcp_congestion_control=bbr' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
</code></pre>
现在，您可以使用以下命令确认BBR已启用：
<pre><code>sudo sysctl net.ipv4.tcp_available_congestion_control
</code></pre>
输出应该类似于：
<pre><code>net.ipv4.tcp_available_congestion_control = bbr cubic reno
</code></pre>
接下来，验证：
<pre><code>sudo sysctl -n net.ipv4.tcp_congestion_control
</code></pre>
输出应该是：
<pre><code>bbr
</code></pre>
最后，检查内核模块是否已加载：
<pre><code>lsmod | grep bbr
</code></pre>
输出将类似于：
<pre><code>tcp_bbr                16384  0
</code></pre>
<h3>步骤3（可选）：测试网络性能增强</h3>
为了测试BBR的网络性能增强，您可以在Web服务器目录中创建一个文件进行下载，然后从台式机上的Web浏览器测试下载速度。
<pre><code>sudo yum install httpd -y
sudo systemctl start httpd.service
sudo firewall-cmd --zone=public --permanent --add-service=http
sudo firewall-cmd --reload
cd /var/www/html
sudo dd if=/dev/zero of=500mb.zip bs=1024k count=500
</code></pre>
最后，<code>http://[your-server-IP]/500mb.zip</code>从桌面计算机上的网络浏览器访问URL ，然后评估下载速度。

就这样。谢谢你的阅读。