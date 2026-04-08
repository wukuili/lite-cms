---
title: "关于 Svelte UI 20：基于 Svelte.js 的轻量级 UI 组件库"
categories: [ "日常" ]
tags: [ "Svelte UI 20","Svelte.js","轻量级UI组件库","无虚拟DOM开销","响应式设计" ]
draft: false
slug: "guan-yu-sveltesvelteui20-ji-yu-sveltejs-qing-liang-ji-ui-zu-jian-ku"
date: "2025-12-05 09:35:53"
url: "/guan-yu-sveltesvelteui20-ji-yu-sveltejs-qing-liang-ji-ui-zu-jian-ku.html"
---

Svelte UI 20 是一款专为 Svelte/SvelteKit 打造的轻量级 UI 组件库，核心特点是**无虚拟 DOM 开销**、**体积小**、**性能优**，贴合 Svelte 「写更少代码，做更多事」的设计理念，适合构建轻量、高效的前端应用。以下从核心特性、快速上手、核心组件示例、定制化等维度详细介绍：

---

## 一、核心特性
1. **极致轻量化**：组件按需引入，无冗余代码，整体打包体积 < 50KB（gzip）；
2. **原生 Svelte 体验**：基于 Svelte 原生语法开发，无额外运行时依赖，组件更新直接编译为原生 DOM 操作；
3. **响应式设计**：适配移动端/PC 端，内置响应式工具类；
4. **主题定制**：支持 CSS 变量/主题文件自定义，适配不同品牌风格；
5. **无障碍支持**：遵循 WAI-ARIA 规范，保障键盘导航、屏幕阅读器兼容性；
6. **TypeScript 友好**：全组件类型声明，开发时类型提示完善。

---

## 二、快速上手
### 1. 安装（基于 SvelteKit 项目）
```bash
# npm
npm install svelte-ui20 --save

# pnpm（推荐）
pnpm add svelte-ui20
```

### 2. 全局引入基础样式（src/routes/+layout.svelte）
```svelte
<script>
  // 引入全局样式（包含重置样式、基础变量）
  import 'svelte-ui20/dist/styles/index.css';
</script>
```

### 3. 按需引入组件
无需额外配置，直接导入使用：
```svelte
<script>
  import { Button, Card, Input } from 'svelte-ui20';

  let value = '';
  const handleClick = () => alert(`输入内容：${value}`);
</script>

<Card title="基础示例" padding="16px">
  <Input 
    v-model={value} 
    placeholder="请输入内容" 
    style="margin-bottom: 12px"
  />
  <Button type="primary" on:click={handleClick}>
    点击提交
  </Button>
</Card>
```

---

## 三、核心组件示例
### 1. 按钮（Button）
支持多种类型、尺寸、状态：
```svelte
<script>
  import { Button } from 'svelte-ui20';
</script>

<Button type="primary">主要按钮</Button>
<Button type="secondary">次要按钮</Button>
<Button type="danger" size="small" disabled>禁用危险按钮</Button>
<Button icon="search" circle>搜索</Button>
```

### 2. 弹窗（Modal）
异步/同步调用，支持自定义内容：
```svelte
<script>
  import { Modal, Button } from 'svelte-ui20';
  let showModal = false;

  async function openModal() {
    const result = await Modal.confirm({
      title: '确认操作',
      content: '是否删除该条数据？',
      okText: '确认',
      cancelText: '取消'
    });
    if (result) {
      console.log('用户确认删除');
    }
  }
</script>

<Button on:click={() => showModal = true}>打开自定义弹窗</Button>
<Button on:click={openModal}>打开确认弹窗</Button>

<Modal 
  bind:visible={showModal}
  title="自定义弹窗"
  on:close={() => showModal = false}
>
  <p>这是自定义弹窗内容</p>
  <div slot="footer">
    <Button on:click={() => showModal = false}>取消</Button>
    <Button type="primary" on:click={() => showModal = false}>确认</Button>
  </div>
</Modal>
```

### 3. 表格（Table）
支持排序、筛选、分页：
```svelte
<script>
  import { Table, Pagination } from 'svelte-ui20';

  // 模拟数据
  const columns = [
    { key: 'name', title: '姓名', sortable: true },
    { key: 'age', title: '年龄', sortable: true },
    { key: 'address', title: '地址' },
    { key: 'action', title: '操作', render: (row) => 
      <Button size="small" type="primary">编辑</Button>
    }
  ];

  const data = [
    { name: '张三', age: 20, address: '北京' },
    { name: '李四', age: 25, address: '上海' },
    { name: '王五', age: 30, address: '广州' }
  ];

  let currentPage = 1;
  const pageSize = 10;
</script>

<Table 
  columns={columns} 
  data={data} 
  bordered
  striped
/>
<Pagination 
  bind:currentPage 
  total={100} 
  pageSize={pageSize}
  style="margin-top: 16px; text-align: right"
/>
```

---

## 四、主题定制
### 1. 覆盖 CSS 变量（全局）
在 `src/app.css` 中定义：
```css
:root {
  /* 主色调 */
  --sui-color-primary: #409eff;
  /* 成功色 */
  --sui-color-success: #67c23a;
  /* 按钮圆角 */
  --sui-btn-border-radius: 4px;
  /* 字体大小 */
  --sui-font-size-base: 14px;
}
```

### 2. 自定义主题文件
创建 `theme-custom.css`，引入并覆盖默认样式：
```css
/* 引入默认主题基础 */
@import 'svelte-ui20/dist/styles/base.css';

/* 自定义组件样式 */
.sui-button--primary {
  background: #722ed1;
  border-color: #722ed1;
}
.sui-button--primary:hover {
  background: #8647e8;
  border-color: #8647e8;
}
```

在 `+layout.svelte` 中引入：
```svelte
<script>
  import './theme-custom.css';
  // 不再引入默认样式
  // import 'svelte-ui20/dist/styles/index.css';
</script>
```

---

## 五、进阶技巧
### 1. 组件懒加载（SvelteKit）
利用 SvelteKit 的动态导入，减少首屏体积：
```svelte
<script>
  let Modal;
  let showModal = false;

  // 点击时加载组件
  async function loadModal() {
    ({ Modal } = await import('svelte-ui20'));
    showModal = true;
  }
</script>

<Button on:click={loadModal}>懒加载弹窗</Button>

{#if Modal}
  <Modal bind:visible={showModal} title="懒加载示例">
    弹窗内容
  </Modal>
{/if}
```

### 2. 自定义组件扩展
基于现有组件封装业务组件：
```svelte
<!-- src/lib/MyButton.svelte -->
<script>
  import { Button } from 'svelte-ui20';
  export let type = 'primary';
  export let size = 'medium';
</script>

<Button 
  {type} 
  {size} 
  class="my-button"
  {...$$restProps}
>
  <slot />
</Button>

<style>
  .my-button {
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }
</style>
```

---

## 六、注意事项
1. **适配 Svelte 版本**：Svelte UI 20 要求 Svelte ≥ 4.0 或 SvelteKit ≥ 1.0；
2. **服务端渲染（SSR）**：部分组件（如 Modal、Tooltip）需在客户端挂载，可通过 `browser` 检测：
   ```svelte
   <script>
     import { browser } from '$app/environment';
     import { Modal } from 'svelte-ui20';
   </script>

   {#if browser}
     <Modal />
   {/if}
   ```
3. **按需引入优化**：配合 `vite-plugin-svelte-kit` 可进一步优化打包体积。

---

## 总结
Svelte UI 20 是 Svelte 生态中轻量化、高性能的 UI 组件库，核心优势是贴合 Svelte 原生语法，无额外运行时开销，适合追求极致性能的中小型应用、管理后台、H5 页面等场景。如果需要更复杂的交互（如可视化、富文本），可结合 Svelte 生态其他库（如 `svelte-chartjs`、`svelte-quill`）扩展。