---
title: "被客户反复退订的cursor又强力了"
categories: [ "未分类" ]
tags: [  ]
draft: true
slug: "trashed"
date: "2025-10-30 10:11:15"
url: "/trashed.html"
---

<!-- wp:heading {"level":3} -->
<h3 class="wp-block-heading">一、工具增强</h3>
<!-- /wp:heading -->

<!-- wp:list {"ordered":true} -->
<ol class="wp-block-list"><!-- wp:list-item -->
<li>1.​<strong>​Browser Controls​</strong>​：类似 Playwright MCP，支持在独立窗口或编辑器内联面板中操作浏览器，可直接选中网页元素同步代码到对话框，实现精准 UI 调整。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_510ce1f3cb6345e69dde8251b6893a22.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789572%3B2077149572&amp;q-key-time=1761789572%3B2077149572&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=69356362ee0f68dc0b93f8d9a515a63656b3702e" alt=""></li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>2.​<strong>​内置工具优化​</strong>​：<!-- wp:list -->
<ul class="wp-block-list"><!-- wp:list-item -->
<li>•​<strong>​Read file​</strong>​：取消文件大小限制，支持完整读取大型文件（如 CSV/JSON）；v1.7 新增图像文件读取能力。</li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>•​<strong>​List/Grep/Codebase Search​</strong>​：提升目录遍历效率、搜索结果相关性及代码索引准确性。</li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>•​<strong>​Web Search​</strong>​：改用轻量级模型，响应更简洁。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_86f3b9db5309440c99f3647b36460698.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789574%3B2077149574&amp;q-key-time=1761789574%3B2077149574&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=d7b4b1c8908ad14f10307771bb1fda4c48433cea" alt=""></li>
<!-- /wp:list-item --></ul>
<!-- /wp:list --></li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>3.​<strong>​Sandboxed terminals​</strong>​：白名单模式下，危险命令自动在沙盒中运行，保障本地系统安全。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_1af729fad1ff4d1d865d80edd9b3780a.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789575%3B2077149575&amp;q-key-time=1761789575%3B2077149575&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=2c4c6e104f1e9b66424a807067fc1d38d35450f2" alt=""></li>
<!-- /wp:list-item --></ol>
<!-- /wp:list -->

<!-- wp:heading {"level":3} -->
<h3 class="wp-block-heading">二、上下文工程改进</h3>
<!-- /wp:heading -->

<!-- wp:list {"ordered":true} -->
<ol class="wp-block-list"><!-- wp:list-item -->
<li>1.​<strong>​Plan Mode​</strong>​：通过生成结构化任务清单（plan.md）推进复杂任务，实现规范驱动开发。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_2a8e9aa0e2f14e72a9bf1620a3f2939b.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789577%3B2077149577&amp;q-key-time=1761789577%3B2077149577&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=1d9f2b4a8ecc85fe04f10ab99639fdb64a56f928" alt=""></li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>2.​<strong>​Hooks​</strong>​：设置自动化规则（如阻止危险命令），在特定事件节点强制触发，提升工程管控效率。</li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>3.​<strong>​Slash Commands​</strong>​：封装常用工作流（如代码审查）为一键命令，支持自定义与快速调用。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_5395f69c81414227b76d2346c6677765.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789579%3B2077149579&amp;q-key-time=1761789579%3B2077149579&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=ce8ccdd90de0422e0fce7e5b346df9d02a696b43" alt=""></li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>4.​<strong>​Autocomplete for Agent​</strong>​：在输入 Prompt 时基于上下文提供补全建议（如变量名）。</li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>5.​<strong>​上下文使用可视化​</strong>​：对话框中直接查看上下文窗口占用情况，搭配自动总结功能优化资源管理。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_e591a7475f5a458188090151a4fb3202.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789580%3B2077149580&amp;q-key-time=1761789580%3B2077149580&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=3e762f1df011b26c4ce25edcdee3da95de121d35" alt=""></li>
<!-- /wp:list-item --></ol>
<!-- /wp:list -->

<!-- wp:heading {"level":3} -->
<h3 class="wp-block-heading">三、v2.0 重要更新</h3>
<!-- /wp:heading -->

<!-- wp:list {"ordered":true} -->
<ol class="wp-block-list"><!-- wp:list-item -->
<li>1.​<strong>​Agent/Editor 布局切换​</strong>​：<!-- wp:list -->
<ul class="wp-block-list"><!-- wp:list-item -->
<li>•​<strong>​Editor 模式​</strong>​：标准 IDE 界面。</li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>•​<strong>​Agent 模式​</strong>​：极简布局，聚焦对话与任务切换，结合 GUI 可视化与 CLI 沉浸感。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_42b55f49ae02426d838c3bbc7b856355.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789583%3B2077149583&amp;q-key-time=1761789583%3B2077149583&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=31eeda2deeeab36518e3b56c674856e356a3f1b8" alt=""></li>
<!-- /wp:list-item --></ul>
<!-- /wp:list --></li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>2.​<strong>​语音输入​</strong>​：内置语音功能，支持更自然的 Vibe Coding。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_6b45f05aa16646498bb845bfbbb2a112.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789584%3B2077149584&amp;q-key-time=1761789584%3B2077149584&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=d4e51bc4bed06d2e25c094e11054c5bb5aac91e6" alt=""></li>
<!-- /wp:list-item -->

<!-- wp:list-item -->
<li>3.​<strong>​Worktrees​</strong>​：基于 Git 功能实现多分支并行，允许同时运行多个 Agent 或测试不同模型建议。<img src="https://hunyuan-plugin-private-1258344706.cos.ap-nanjing.myqcloud.com/pdf_youtu/img/1761789567_f734551502c4453ea20e5a507cd4a966.png?q-sign-algorithm=sha1&amp;q-ak=AKID372nLgqocp7HZjfQzNcyGOMTN3Xp6FEA&amp;q-sign-time=1761789585%3B2077149585&amp;q-key-time=1761789585%3B2077149585&amp;q-header-list=host&amp;q-url-param-list=&amp;q-signature=3d78a5f9219521915f258b7e11219977daef54ac" alt=""></li>
<!-- /wp:list-item --></ol>
<!-- /wp:list -->

<!-- wp:heading {"level":3} -->
<h3 class="wp-block-heading">四、总结</h3>
<!-- /wp:heading -->

<!-- wp:paragraph -->
<p>Cursor 通过高频迭代在工具链、上下文工程及用户体验上大幅提升，尤其在安全管控（沙盒）、自动化（Hooks）与多任务处理（Worktrees）方面形成差异化优势，重新巩固了其在 AI 编程工具中的竞争力。</p>
<!-- /wp:paragraph -->