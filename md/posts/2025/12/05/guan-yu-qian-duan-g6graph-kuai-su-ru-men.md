---
title: "前端 G6 图可视化库快速入门"
categories: [ "日常" ]
tags: [ "G6 图可视化库","关系型数据可视化","力导向图","G6 快速入门","前端可视化" ]
draft: false
slug: "guan-yu-qian-duan-g6graph-kuai-su-ru-men"
date: "2025-12-05 09:41:43"
url: "/guan-yu-qian-duan-g6graph-kuai-su-ru-men.html"
---

G6 是蚂蚁集团开源的图可视化引擎，专注于关系型数据的可视化展示与交互，支持流程图、拓扑图、知识图谱、力导向图等多种场景。以下是**从环境搭建到核心案例**的快速上手指南，覆盖基础使用、核心配置、交互扩展等关键内容。

---

## 一、核心特性
- 高性能：支持万级节点/边的渲染与交互（Canvas 渲染核心）；
- 丰富的内置元素：节点、边、锚点等基础元素，支持自定义形状；
- 灵活的布局：力导向、树形、环形、分层（Dagre）等 20+ 内置布局；
- 完整的交互体系：拖拽、缩放、选中等基础交互，支持自定义交互行为；
- 插件化扩展：内置 tooltip、mini-map、undo/redo 等插件，支持自定义插件。

---

## 二、快速上手（基础环境）
### 1. 安装
#### 方式 1：NPM 安装（推荐，适用于 Vue/React/Svelte 等工程化项目）
```bash
# 最新稳定版
npm install @antv/g6 --save

# 或 yarn/pnpm
yarn add @antv/g6
pnpm add @antv/g6
```

#### 方式 2：CDN 引入（适用于纯 HTML 项目）
```html
<!-- 引入 G6 核心库 -->
<script src="https://gw.alipayobjects.com/os/lib/antv/g6/5.0.0/dist/g6.min.js"></script>
```

### 2. 最小示例（渲染基础力导向图）
```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>G6 快速入门</title>
  <!-- 容器必须设置宽高 -->
  <style>
    #container {
      width: 800px;
      height: 600px;
      border: 1px solid #eee;
    }
  </style>
</head>
<body>
  <div id="container"></div>

  <script>
    // 1. 准备数据（节点 + 边）
    const data = {
      // 节点数组：每个节点必须有 id，name/shape/style 为可选
      nodes: [
        { id: 'node1', name: '节点1', x: 100, y: 200 },
        { id: 'node2', name: '节点2', x: 300, y: 200 },
        { id: 'node3', name: '节点3', x: 200, y: 100 },
      ],
      // 边数组：每个边必须有 source（源节点）、target（目标节点）
      edges: [
        { source: 'node1', target: 'node2' },
        { source: 'node1', target: 'node3' },
        { source: 'node2', target: 'node3' },
      ],
    };

    // 2. 初始化图实例
    const graph = new G6.Graph({
      container: 'container', // 挂载容器 ID
      width: 800, // 图宽度
      height: 600, // 图高度
      // 基础配置
      defaultNode: { // 节点默认样式
        size: 60, // 节点大小
        style: { fill: '#409eff', stroke: '#fff', lineWidth: 2 },
        labelCfg: { style: { fill: '#fff', fontSize: 14 } }, // 节点文字样式
      },
      defaultEdge: { // 边默认样式
        style: { stroke: '#999', lineWidth: 2 },
      },
      // 布局：力导向布局（无坐标时自动计算节点位置）
      layout: {
        type: 'force', // 布局类型
        linkDistance: 150, // 边的长度
        preventOverlap: true, // 防止节点重叠
      },
    });

    // 3. 加载数据并渲染
    graph.data(data);
    graph.render();

    // 4. 开启交互：缩放、拖拽
    graph.enableZoom();
    graph.enablePan();
  </script>
</body>
</html>
```

#### 效果说明：
- 3 个节点自动以力导向布局排列，节点间有边连接；
- 支持鼠标滚轮缩放、拖拽画布平移、拖拽节点调整位置。

---

## 三、核心概念与配置
### 1. 数据格式（核心）
G6 严格要求数据为「节点 + 边」的结构化格式：
```javascript
const data = {
  nodes: [
    {
      id: '唯一标识', // 必选
      // 可选配置
      name: '显示文本',
      shape: 'circle', // 节点形状：circle/rect/ellipse 等
      size: 50, // 节点大小
      x: 100, y: 200, // 固定坐标（布局为 none 时生效）
      style: { fill: '#f00' }, // 自定义样式
      labelCfg: { position: 'bottom' }, // 文字位置
    },
  ],
  edges: [
    {
      source: '源节点 ID', // 必选
      target: '目标节点 ID', // 必选
      // 可选配置
      shape: 'line', // 边形状：line/curve/loop 等
      label: '边文字', // 边的标注
      style: { stroke: '#00f', lineWidth: 3 },
    },
  ],
};
```

### 2. 布局配置（控制节点排列）
G6 内置多种布局，适配不同场景：
| 布局类型 | 适用场景 | 核心配置 |
|----------|----------|----------|
| `force`（力导向） | 无层级的关系图（如社交网络） | `linkDistance`（边长度）、`preventOverlap`（防重叠） |
| `tree`（树形） | 层级结构（如组织架构） | `direction`（布局方向：TB/BT/LR/RL）、`depthFactor`（层级间距） |
| `dagre`（分层） | 流程图、DAG 图 | `rankdir`（方向）、`nodesep`（节点间距） |
| `circular`（环形） | 环形关系图 | `radius`（半径）、`startAngle`（起始角度） |
| `none`（无布局） | 自定义坐标 | 需手动设置节点 x/y |

#### 树形布局示例：
```javascript
layout: {
  type: 'tree',
  direction: 'TB', // 从上到下布局
  nodeSep: 50, // 同层节点间距
  rankSep: 100, // 层级间距
}
```

### 3. 自定义节点/边
#### 示例：自定义矩形节点 + 带箭头的边
```javascript
const graph = new G6.Graph({
  container: 'container',
  width: 800,
  height: 600,
  defaultNode: {
    shape: 'rect', // 矩形节点
    size: [120, 60], // 矩形宽高（数组）
    style: { fill: '#67c23a', radius: 8 }, // 圆角矩形
    labelCfg: { style: { fontSize: 16 } },
  },
  defaultEdge: {
    shape: 'polyline', // 折线边
    style: { 
      stroke: '#409eff', 
      lineWidth: 2,
      endArrow: true, // 显示终点箭头
    },
  },
  layout: { type: 'force' },
});
```

### 4. 交互配置
G6 内置丰富的交互能力，通过简单 API 开启：
```javascript
// 基础交互
graph.enableZoom(); // 缩放（滚轮）
graph.enablePan(); // 平移（拖拽画布）
graph.enableDragNode(); // 拖拽节点

// 高级交互：节点点击/悬停
graph.on('node:click', (e) => {
  const node = e.item; // 获取点击的节点
  console.log('点击节点：', node.getModel()); // 获取节点数据
  // 高亮点击的节点
  graph.setItemState(node, 'active', true);
});

// 取消节点高亮（点击空白处）
graph.on('canvas:click', () => {
  graph.clearItemStates();
});

// 边点击事件
graph.on('edge:click', (e) => {
  console.log('点击边：', e.item.getModel());
});
```

---

## 四、常用插件（提升体验）
### 1. 提示框（Tooltip）
```javascript
// 引入插件（NPM 方式需导入，CDN 已内置）
import { Tooltip } from '@antv/g6';

// 注册插件
graph.registerPlugin(
  new Tooltip({
    offset: 10, // 提示框偏移量
    // 自定义提示框内容
    getContent: (e) => {
      const model = e.item.getModel();
      if (e.item.getType() === 'node') {
        return `<div style="padding: 8px;">节点：${model.name}</div>`;
      }
      return `<div style="padding: 8px;">边：${model.source} → ${model.target}</div>`;
    },
  })
);
```

### 2. 迷你地图（MiniMap）
```javascript
import { MiniMap } from '@antv/g6';

graph.registerPlugin(
  new MiniMap({
    size: [150, 100], // 迷你地图大小
    className: 'minimap', // 自定义类名
    style: { fill: '#f8f8f8' }, // 背景样式
  })
);
```

### 3. 撤销/重做（UndoRedo）
```javascript
import { UndoRedo } from '@antv/g6';

const undoRedo = new UndoRedo({ graph });
graph.registerPlugin(undoRedo);

// 绑定快捷键
graph.on('keydown', (e) => {
  if (e.key === 'z' && e.ctrlKey) undoRedo.undo(); // Ctrl+Z 撤销
  if (e.key === 'y' && e.ctrlKey) undoRedo.redo(); // Ctrl+Y 重做
});
```

---

## 五、实战案例：组织架构图（树形布局）
```javascript
// 组织架构数据
const orgData = {
  nodes: [
    { id: '1', name: 'CEO', size: 80, style: { fill: '#e6a23c' } },
    { id: '2', name: '产品部', size: 60 },
    { id: '3', name: '研发部', size: 60 },
    { id: '4', name: '设计部', size: 60 },
    { id: '5', name: '前端组', size: 50 },
    { id: '6', name: '后端组', size: 50 },
  ],
  edges: [
    { source: '1', target: '2' },
    { source: '1', target: '3' },
    { source: '1', target: '4' },
    { source: '3', target: '5' },
    { source: '3', target: '6' },
  ],
};

// 初始化图
const graph = new G6.Graph({
  container: 'container',
  width: 800,
  height: 600,
  defaultNode: {
    shape: 'rect',
    size: [100, 60],
    style: { fill: '#409eff', radius: 4 },
    labelCfg: { style: { fill: '#fff', fontSize: 14 } },
  },
  defaultEdge: {
    shape: 'polyline',
    style: { stroke: '#999', lineWidth: 2, endArrow: true },
  },
  // 树形布局（从上到下）
  layout: {
    type: 'tree',
    direction: 'TB',
    nodeSep: 30,
    rankSep: 80,
    center: [400, 50], // 布局中心
  },
});

// 加载数据并渲染
graph.data(orgData);
graph.render();

// 开启交互
graph.enableZoom();
graph.enableDragNode();
```

---

## 六、常见问题与优化
### 1. 大量节点渲染卡顿
- 开启「节点聚合」：使用 `combo` 将关联节点聚合，减少渲染数量；
- 启用「视口裁剪」：只渲染可视区域内的节点/边；
- 切换渲染模式：G6 5.x 支持 WebGL 渲染（`renderer: 'webgl'`），大幅提升万级节点性能。

### 2. 节点/边样式不生效
- 检查 `style` 配置是否符合 G6 规范（如 `fill` 对应填充，`stroke` 对应描边）；
- 自定义样式优先级高于 `defaultNode/defaultEdge`，需注意覆盖顺序。

### 3. 布局不生效
- 若手动设置了节点 `x/y`，需将 `layout.type` 设为 `none`；
- 力导向布局需等待计算完成，可通过 `graph.on('layoutend', () => {})` 监听布局完成。

---

## 七、进阶学习资源
1. **官方文档**：[G6 官方文档](https://g6.antv.vision/zh)（包含完整 API、案例、原理）；
2. **示例库**：[G6 示例中心](https://g6.antv.vision/zh/examples)（流程图、知识图谱、大屏可视化等）；
3. **源码解析**：G6 核心源码在 GitHub（antvis/G6），可学习 Canvas 渲染、布局算法实现。

---

## 总结
G6 上手核心是「理解数据格式 + 配置布局 + 自定义样式/交互」：
- 基础场景（如简单关系图）仅需 5 分钟即可实现；
- 复杂场景（如组织架构、海量节点）可通过布局、插件、性能优化扩展；
- 适配 Vue/React 项目时，只需将 G6 实例挂载到组件生命周期（如 Vue 的 `onMounted` 初始化，`onUnmounted` 销毁）。