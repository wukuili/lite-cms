---
title: "轻松去除 Alist 端口号"
categories: [ "日常" ]
tags: [ "Alist端口号去除","反向代理配置","Nginx配置","HTTPS加密访问","域名解析" ]
draft: false
slug: "qing-song-qu-chu-alist-duan-kou-hao-jian-dan-bu-zhou-ti-sheng-nin-de-wang-luo-ti-yan"
date: "2025-11-23 21:16:04"
url: "/qing-song-qu-chu-alist-duan-kou-hao-jian-dan-bu-zhou-ti-sheng-nin-de-wang-luo-ti-yan.html"
---

轻松去除Alist端口号，核心是通过反向代理配置实现，无需复杂操作，能显著优化网络访问体验。

### 一、核心原理
1. 利用Nginx或Caddy等反向代理工具，将域名直接映射到Alist服务的“IP:端口”。
2. 访问时只需输入域名（如alist.example.com），代理工具会自动转发请求到对应端口，实现端口号隐藏。

### 二、简单操作步骤（以Nginx为例）
1. 安装Nginx：通过服务器包管理工具快速安装（如CentOS用`yum install nginx`，Ubuntu用`apt install nginx`）。
2. 配置反向代理：编辑Nginx配置文件（通常在`/etc/nginx/conf.d/`目录下），添加如下内容：
   ```
   server {
       listen 80;
       server_name alist.example.com; # 替换为你的域名
       location / {
           proxy_pass http://127.0.0.1:5244; # 替换为Alist实际的IP:端口
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```
3. 重启Nginx：执行命令`systemctl restart nginx`，使配置生效。
4. 解析域名：在域名服务商处，将域名解析到服务器IP，等待解析生效后即可通过域名访问Alist。

### 三、注意事项
1. 确保服务器80端口（HTTP）或443端口（HTTPS）未被占用，且已在防火墙放行。
2. 若需HTTPS加密访问，可通过Let’s Encrypt申请免费证书，在Nginx配置中添加SSL相关参数。

---

以下是 **Alist 端口号隐藏的 Nginx 完整配置文件**（含 HTTP 自动跳转 HTTPS + 免费 SSL 证书配置），直接复制替换参数即可使用，无需额外修改核心逻辑：

### 一、完整配置文件（复制即用）
```nginx
# 保存路径：/etc/nginx/conf.d/alist.conf（新建文件，无需修改主配置）
server {
    # 监听 80 端口（HTTP），自动跳转到 HTTPS
    listen 80;
    listen [::]:80;
    server_name alist.example.com; # 替换为你的域名（如 alist.你的域名.com）

    # 强制 HTTP 转 HTTPS（无需手动输入 https，更友好）
    return 301 https://$host$request_uri;
}

server {
    # 监听 443 端口（HTTPS），隐藏端口号核心配置
    listen 443 ssl;
    listen [::]:443 ssl;
    server_name alist.example.com; # 同上，保持域名一致

    # 1. SSL 证书配置（Let's Encrypt 免费证书，自动生成路径）
    ssl_certificate /etc/letsencrypt/live/alist.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/alist.example.com/privkey.pem;
    ssl_trusted_certificate /etc/letsencrypt/live/alist.example.com/chain.pem;

    # 2. SSL 优化配置（提升安全性和兼容性，无需修改）
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers "EECDH+AESGCM:EDH+AESGCM:AES256+EECDH:AES256+EDH";
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;

    # 3. 反向代理核心配置（映射 Alist 端口）
    location / {
        proxy_pass http://127.0.0.1:5244; # 替换为你的 Alist 实际地址（默认 5244 端口）
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme; # 传递 HTTPS 协议标识，避免 Alist 识别异常
        proxy_buffering off; # 关闭缓冲，提升文件上传下载速度
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 4. 静态资源缓存（优化访问速度，可选保留）
    location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
        proxy_pass http://127.0.0.1:5244;
        proxy_set_header Host $host;
        expires 30d; # 静态资源缓存 30 天
        add_header Cache-Control "public, max-age=2592000";
    }
}
```

### 二、配置步骤（3 分钟完成）
#### 1. 准备工作
- 已解析域名：将 `alist.example.com` 解析到你的服务器 IP（域名服务商控制台操作，解析类型 A 记录）。
- 服务器放行端口：防火墙开放 80（HTTP）、443（HTTPS）端口（如阿里云/腾讯云安全组添加规则）。

#### 2. 安装免费 SSL 证书（Let's Encrypt）
执行以下命令自动申请证书（需先安装 Certbot）：
```bash
# 1. 安装 Certbot（Ubuntu/Debian 示例，CentOS 用 yum 替换 apt）
sudo apt update && sudo apt install certbot python3-certbot-nginx -y

# 2. 申请证书（自动配置 Nginx，无需手动复制证书文件）
sudo certbot --nginx -d alist.example.com
```
- 执行后按提示操作：输入邮箱 → 同意协议 → 选择是否共享邮箱 → 自动配置 SSL（选 1 或 2 均可）。

#### 3. 部署配置文件
```bash
# 1. 新建 Nginx 配置文件
sudo nano /etc/nginx/conf.d/alist.conf

# 2. 复制上面的完整配置，粘贴到文件中，修改 2 个关键参数：
#    - server_name：替换为你的域名
#    - proxy_pass：替换为你的 Alist 地址（默认 127.0.0.1:5244，若 Alist 部署在其他服务器则改对应 IP）

# 3. 保存退出（nano 编辑器：按 Ctrl+O → 回车 → Ctrl+X）

# 4. 检查配置是否有误
sudo nginx -t

# 5. 重启 Nginx 生效
sudo systemctl restart nginx
```

### 三、验证效果
1. 访问 `https://你的域名`（无需加端口号），直接进入 Alist 登录页 → 配置成功。
2. 若之前能通过 `IP:5244` 访问，现在仍可正常使用（端口未关闭，仅新增域名访问方式）。

### 四、常见问题解决
1. **Nginx 启动失败**：检查配置文件语法（`nginx -t` 会提示错误位置），确保域名、Alist 端口填写正确。
2. **访问域名跳转到服务器默认页**：可能是 Nginx 存在其他默认配置（如 `default.conf`），可删除或重命名：`sudo mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.bak`，再重启 Nginx。
3. **Alist 上传/下载异常**：确保 `proxy_buffering off` 和 `proxy_http_version 1.1` 配置已添加，避免缓冲导致的断点续传问题。
4. **证书自动续期**：Let's Encrypt 证书有效期 90 天，执行以下命令设置自动续期：
   ```bash
   sudo crontab -e
   # 粘贴以下内容（每天凌晨 2 点自动检查续期）
   0 2 * * * /usr/bin/certbot renew --quiet
   ```

按以上配置，你可以直接通过 `https://你的域名` 访问 Alist，完全隐藏 5244 端口，同时获得 HTTPS 加密的安全访问体验～