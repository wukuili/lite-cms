---
title: "关于JavaScript你不知道的defaultprevented"
categories: [ "日常" ]
tags: [ "JavaScript","defaultPrevented","事件对象","preventDefault","事件默认行为" ]
draft: false
slug: "guan-yu-javascript-ni-bu-zhi-dao-de-preventdefault"
date: "2025-11-19 13:58:25"
url: "/guan-yu-javascript-ni-bu-zhi-dao-de-preventdefault.html"
---

`defaultPrevented` 是 JavaScript 事件对象（Event）的只读属性，核心结论是：用于判断事件的默认行为是否已通过 `event.preventDefault()` 被取消。

### 核心说明
- 类型：布尔值（true/false）。
- 触发逻辑：调用 `preventDefault()` 后，`defaultPrevented` 会变为 true；未调用或事件无法取消默认行为时，值为 false。
- 用途：区分事件默认行为是“未触发”还是“被主动取消”，避免误判。

### 简单示例
```javascript
// 阻止链接跳转（默认行为）
document.querySelector('a').addEventListener('click', (e) => {
  e.preventDefault(); // 取消默认跳转行为
  console.log(e.defaultPrevented); // 输出 true
});

// 未阻止默认行为的事件
document.addEventListener('scroll', (e) => {
  console.log(e.defaultPrevented); // 输出 false（scroll 事件默认行为无法取消）
});
```
