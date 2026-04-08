---
title: "关于 Wing FTP Server 搭建 FTP 服务器的入门教程"
categories: [ "日常" ]
tags: [ "Wing FTP Server","FTP服务器搭建","Windows系统使用","FTP连接测试","新手教程" ]
draft: false
slug: "guan-yu-ftpwingftpserver-da-jian-ftp-fu-wu-qi-ru-men-jiao-cheng"
date: "2025-11-24 15:12:43"
url: "/guan-yu-ftpwingftpserver-da-jian-ftp-fu-wu-qi-ru-men-jiao-cheng.html"
---

### 一、核心结论
Wing FTP Server 是一款轻量易操作的 FTP 服务器软件，Windows/macOS/Linux 均支持，新手通过“下载安装→创建域名→添加用户→配置权限→测试连接”5步即可完成基础搭建。

---

### 二、关键步骤（以 Windows 系统为例）
#### 1. 下载与安装软件
- 官网下载最新版 Wing FTP Server（有免费版，满足个人/小型团队使用）。
- 安装时默认下一步即可，勾选“创建桌面快捷方式”方便后续操作。

#### 2. 创建 FTP 域名（核心配置）
- 打开软件，左侧导航栏点击“Domains”→“Add Domain”。
- 填写域名名称（自定义，如“MyFTP”），端口默认 21（FTP 标准端口，无冲突建议保留）。
- 根目录选择本地文件夹（如 D:\FTPFile），作为 FTP 服务器的文件存储路径，点击“OK”保存。

#### 3. 添加 FTP 访问用户
- 选中已创建的域名→右侧“Users”→“Add User”。
- 填写用户名（如“ftpuser”）和密码，勾选“Enable Account”激活账号。
- 切换到“Dir Access”标签，选择根目录，设置权限（如“Read+Write”读写权限，根据需求调整）。

#### 4. 关闭防火墙拦截（重要）
- 进入 Windows 防火墙设置，允许 Wing FTP Server 程序通过防火墙（包括公用/专用网络）。
- 若需外网访问，需在路由器中做“端口映射”（映射本地 IP + 21 端口）。

#### 5. 测试 FTP 连接
- 本地测试：使用 FileZilla 等 FTP 客户端，输入服务器 IP（本地为 127.0.0.1）、用户名、密码，端口 21，连接成功即可上传/下载文件。
- 外网测试：用外网设备输入服务器公网 IP，重复上述连接步骤，成功则搭建完成。

---

### 三、新手注意事项
- 免费版支持最多 3 个并发连接，商用需升级授权。
- 避免使用默认端口 21 暴露在外网（可修改为 1024+ 端口，提升安全性）。
- 根目录建议单独创建文件夹，不要直接用系统盘（如 C 盘），防止误操作导致系统文件风险。

---
