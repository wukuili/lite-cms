---
title: "Vue 项目作为子应用接入乾坤（qiankun）框架改造笔记"
categories: [ "日常" ]
tags: [ "Vue项目","qiankun框架","子应用接入","改造笔记","前端开发" ]
draft: false
slug: "guan-yu-qiankunvue-xiang-mu-zuo-wei-zi-ying-yong-jie-ru-qiankun-kuang-jia-gai-zao-bi-ji"
date: "2025-11-23 21:28:12"
url: "/guan-yu-qiankunvue-xiang-mu-zuo-wei-zi-ying-yong-jie-ru-qiankun-kuang-jia-gai-zao-bi-ji.html"
---

## 一、改造核心目标
将独立的 Vue 项目（Vue2/Vue3 均适用）改造为 qiankun 可识别的子应用，实现：
1. 主应用通过域名/路径无感知嵌入子应用
2. 主-子应用间通信、路由隔离
3. 样式隔离、资源加载不冲突
4. 支持独立运行（改造后仍可单独启动调试）

## 二、前置准备
1. 环境要求：
   - 主应用已集成 qiankun（参考 [qiankun 官方文档](https://qiankun.umijs.org/zh/guide)）
   - 子应用 Vue 版本：Vue2（需配合 vue-router 3.x）、Vue3（需配合 vue-router 4.x）
   - 打包工具：Webpack（Vue CLI 默认）或 Vite（需额外配置）
2. 关键概念：
   - 子应用需暴露 `bootstrap`、`mount`、`unmount` 三个生命周期钩子
   - 路由需配置为**基于主应用分配的基础路径**（如主应用通过 `/vue-app` 访问子应用，则子应用路由 base 为 `/vue-app`）
   - 资源需支持**跨域访问**（主应用加载子应用静态资源时触发 CORS）

## 三、分步骤改造（Vue CLI 项目，Vue2/Vue3 通用）
### 步骤 1：安装 qiankun 依赖（子应用）
子应用需引入 qiankun 的辅助依赖，用于暴露生命周期：
```bash
# Vue2/Vue3 均适用
npm install qiankun --save
# 或 yarn add qiankun
```

### 步骤 2：配置子应用入口文件（main.js）
核心是暴露 qiankun 要求的三个生命周期钩子，同时兼容独立运行模式。

#### Vue2 示例（main.js）：
```javascript
import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

let instance = null;
let routerInstance = null;

// 初始化 Vue 实例（兼容独立运行和qiankun嵌入）
function render(props = {}) {
  const { container, routerBase = '/' } = props; // 主应用传递的参数（容器、基础路径）
  
  // 初始化路由（基于主应用传递的 base 路径）
  routerInstance = new VueRouter({
    mode: 'history', // 必须使用 history 模式（hash 模式需主应用特殊配置）
    base: window.__POWERED_BY_QIANKUN__ ? routerBase : process.env.BASE_URL, // 独立运行时用默认 base
    routes: router.options.routes // 复用原路由配置
  });

  // 挂载实例（容器优先用主应用传递的 container，独立运行时用 #app）
  instance = new Vue({
    router: routerInstance,
    store,
    render: h => h(App)
  }).$mount(container ? container.querySelector('#app') : '#app');
}

// 独立运行时直接渲染
if (!window.__POWERED_BY_QIANKUN__) {
  render();
}

// 1.  bootstrap：子应用初始化（只执行一次）
export async function bootstrap() {
  console.log('Vue2 子应用 bootstrap');
}

// 2. mount：子应用挂载（主应用切换到子应用时执行）
export async function mount(props) {
  console.log('Vue2 子应用 mount，接收主应用参数：', props);
  render(props); // 传入主应用参数渲染
}

// 3. unmount：子应用卸载（主应用切换离开子应用时执行）
export async function unmount() {
  console.log('Vue2 子应用 unmount');
  instance.$destroy(); // 销毁 Vue 实例
  instance = null;
  routerInstance = null;
}

// 可选：子应用更新（主应用触发子应用更新时执行）
export async function update(props) {
  console.log('Vue2 子应用 update，接收更新参数：', props);
}
```

#### Vue3 示例（main.js）：
```javascript
import { createApp } from 'vue'
import App from './App.vue'
import { createRouter, createWebHistory } from 'vue-router'
import routes from './router/routes' // 抽离的路由配置
import store from './store'

let app = null;
let router = null;

function render(props = {}) {
  const { container, routerBase = '/' } = props;
  
  // 初始化路由
  router = createRouter({
    history: createWebHistory(window.__POWERED_BY_QIANKUN__ ? routerBase : import.meta.env.BASE_URL),
    routes
  });

  // 创建并挂载应用
  app = createApp(App);
  app.use(router);
  app.use(store);
  app.mount(container ? container.querySelector('#app') : '#app');
}

// 独立运行时直接渲染
if (!window.__POWERED_BY_QIANKUN__) {
  render();
}

// 生命周期钩子
export async function bootstrap() {
  console.log('Vue3 子应用 bootstrap');
}

export async function mount(props) {
  console.log('Vue3 子应用 mount，接收主应用参数：', props);
  render(props);
}

export async function unmount() {
  console.log('Vue3 子应用 unmount');
  app.unmount(); // 销毁 Vue3 应用
  app = null;
  router = null;
}

export async function update(props) {
  console.log('Vue3 子应用 update，接收更新参数：', props);
}
```

### 步骤 3：配置路由（router/index.js）
核心要求：路由模式必须为 `history`（hash 模式需主应用配置 `supportHistory: false`，不推荐），且 base 路径动态适配。

#### Vue2 路由示例（router/index.js）：
```javascript
import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '@/views/Home.vue'
import About from '@/views/About.vue'

Vue.use(VueRouter);

// 抽离路由规则（供 main.js 复用）
const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/about', name: 'About', component: About }
];

// 独立运行时的路由实例（嵌入 qiankun 时会重新创建）
const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
});

export default router;
```

#### Vue3 路由示例（router/routes.js）：
```javascript
// 抽离路由规则
export default [
  { path: '/', name: 'Home', component: () => import('@/views/Home.vue') },
  { path: '/about', name: 'About', component: () => import('@/views/About.vue') }
];
```

### 步骤 4：配置打包（vue.config.js）
关键配置：设置 `publicPath`（资源基础路径）、`output`（库模式打包）、`devServer`（跨域支持）。

#### Vue2/Vue3 通用 vue.config.js：
```javascript
const { name } = require('./package.json');

module.exports = {
  // 1. 资源基础路径：嵌入 qiankun 时，主应用需通过绝对路径加载子应用资源
  publicPath: process.env.NODE_ENV === 'production' 
    ? '//你的子应用线上域名/' // 生产环境：子应用部署后的绝对域名（如 https://sub-vue-app.example.com/）
    : '//localhost:8081/', // 开发环境：子应用本地端口（避免和主应用端口冲突）

  // 2. 打包配置：输出为库模式（qiankun 需加载 umd 格式的包）
  configureWebpack: {
    output: {
      library: `${name}-[name]`,
      libraryTarget: 'umd', // 关键：umd 格式支持全局引入和模块化引入
      chunkLoadingGlobal: `webpackJsonp_${name}`, // 避免 chunk 命名冲突
    },
  },

  // 3. 开发环境配置：允许主应用跨域访问（主应用地址需替换为实际地址）
  devServer: {
    port: 8081, // 子应用本地端口（不要和主应用冲突，如主应用 8080，子应用用 8081）
    headers: {
      'Access-Control-Allow-Origin': '*', // 开发环境允许所有跨域（生产环境需指定主应用域名）
      'Access-Control-Allow-Methods': '*',
      'Access-Control-Allow-Headers': '*',
    },
  },
};
```

### 步骤 5：样式隔离配置
避免子应用样式污染主应用，或主应用样式影响子应用，两种方案可选：

#### 方案 1：CSS Modules（推荐）
Vue 组件样式默认支持 `scoped`，但部分全局样式仍可能冲突，建议：
1. 组件内样式统一加 `scoped`：
```vue
<style scoped>
/* 组件内样式，自动添加哈希前缀，避免冲突 */
.container {
  margin: 0;
}
</style>
```
2. 全局样式（如 `main.css`）添加唯一前缀（如 `vue-sub-app-`）：
```css
/* 全局样式前缀 */
.vue-sub-app-container {
  padding: 20px;
}
.vue-sub-app-btn {
  background: #42b983;
}
```

#### 方案 2：qiankun 自带样式隔离
主应用注册子应用时，配置 `sandbox: { strictStyleIsolation: true }`，qiankun 会通过 Shadow DOM 隔离样式：
```javascript
// 主应用注册子应用时的配置
import { registerMicroApps } from 'qiankun';

registerMicroApps([
  {
    name: 'vue-sub-app', // 子应用名称（需和子应用 package.json name 一致）
    entry: '//localhost:8081', // 子应用入口（开发环境本地地址，生产环境线上域名）
    container: '#sub-app-container', // 主应用中挂载子应用的容器
    activeRule: '/vue-app', // 主应用访问子应用的路径（如主应用域名/vue-app 触发加载）
    sandbox: {
      strictStyleIsolation: true, // 开启样式隔离（Vue 子应用推荐开启）
    },
    props: {
      routerBase: '/vue-app', // 传递给子应用的基础路径（需和 activeRule 一致）
      // 其他需要传递的参数（如用户信息、全局配置）
      userInfo: { name: 'test' },
    },
  },
]);
```

### 步骤 6：主-子应用通信配置
qiankun 提供两种通信方式：`props` 传递（一次性初始化参数）和 `initGlobalState`（全局状态管理，支持双向通信）。

#### 方式 1：props 传递（初始化参数）
主应用注册子应用时通过 `props` 传递参数，子应用在 `mount` 钩子中接收：
```javascript
// 主应用注册时传递 props
registerMicroApps([
  {
    name: 'vue-sub-app',
    entry: '//localhost:8081',
    container: '#sub-app-container',
    activeRule: '/vue-app',
    props: {
      routerBase: '/vue-app',
      userInfo: { id: 1, name: '张三' },
      onSubAppMessage: (msg) => { // 主应用接收子应用消息的回调
        console.log('主应用收到子应用消息：', msg);
      },
    },
  },
]);

// 子应用 mount 钩子中接收 props
export async function mount(props) {
  console.log('子应用接收主应用参数：', props.userInfo);
  // 子应用调用主应用回调传递消息
  props.onSubAppMessage('子应用已加载完成');
  render(props);
}
```

#### 方式 2：全局状态管理（双向通信）
适用于主-子应用需要实时同步状态的场景（如用户登录状态、主题切换）：
1. 主应用创建全局状态：
```javascript
// 主应用
import { initGlobalState } from 'qiankun';

// 初始化全局状态
const initialState = {
  theme: 'light',
  userInfo: null,
};
const actions = initGlobalState(initialState);

// 主应用监听状态变化
actions.onGlobalStateChange((state, prev) => {
  console.log('主应用全局状态变化：', state, prev);
});

// 主应用更新状态
actions.setGlobalState({ theme: 'dark' });

// 将 actions 传递给子应用（通过 props）
registerMicroApps([
  {
    name: 'vue-sub-app',
    entry: '//localhost:8081',
    container: '#sub-app-container',
    activeRule: '/vue-app',
    props: {
      routerBase: '/vue-app',
      globalStateActions: actions, // 传递全局状态操作方法
    },
  },
]);
```

2. 子应用使用全局状态：
```javascript
// 子应用 mount 钩子中接收并使用
export async function mount(props) {
  const { globalStateActions } = props;
  
  // 子应用监听全局状态变化
  globalStateActions.onGlobalStateChange((state) => {
    console.log('子应用收到全局状态：', state.theme);
    // 可将状态同步到子应用的 Vuex/Pinia
    store.commit('setTheme', state.theme);
  }, true); // 第二个参数为 true，立即触发一次回调
  
  // 子应用更新全局状态
  globalStateActions.setGlobalState({ userInfo: { id: 1, name: '张三' } });
  
  render(props);
}
```

## 四、测试验证
### 1. 独立运行测试
子应用单独启动，访问 `http://localhost:8081`，确认路由、功能正常。

### 2. 嵌入主应用测试
1. 启动主应用（如 `http://localhost:8080`）
2. 访问主应用的子应用路径（如 `http://localhost:8080/vue-app`）
3. 验证：
   - 子应用正常挂载，无样式冲突
   - 子应用路由正常（如 `/vue-app/about` 访问子应用的 About 页面）
   - 主-子应用通信正常
   - 切换主应用其他页面时，子应用正常卸载

### 3. 生产环境测试
1. 子应用打包：`npm run build`，部署到线上服务器（如 `https://sub-vue-app.example.com`）
2. 主应用注册子应用时，将 `entry` 改为线上域名：
```javascript
registerMicroApps([
  {
    name: 'vue-sub-app',
    entry: 'https://sub-vue-app.example.com', // 生产环境入口
    container: '#sub-app-container',
    activeRule: '/vue-app',
    props: { routerBase: '/vue-app' },
  },
]);
```
3. 访问主应用线上地址，验证子应用加载正常。

## 五、常见问题排查
### 1. 子应用加载失败：`qiankun: application xxx died in status LOADING_SOURCE_CODE`
- 原因：子应用资源跨域（生产环境未配置 CORS）
- 解决：生产环境子应用服务器配置 `Access-Control-Allow-Origin` 为住应用域名（如 Nginx 配置）：
```nginx
server {
  listen 80;
  server_name sub-vue-app.example.com;
  location / {
    root /path/to/vue-sub-app/dist;
    try_files $uri $uri/ /index.html;
    # 允许主应用域名跨域
    add_header 'Access-Control-Allow-Origin' 'https://main-app.example.com';
    add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
    add_header 'Access-Control-Allow-Headers' 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type';
  }
}
```

### 2. 子应用路由跳转 404
- 原因：
  1. 子应用路由 `base` 未配置为主应用传递的 `routerBase`
  2. 主应用未配置子应用路径的 fallback（刷新子应用页面时主应用路由拦截）
- 解决：
  - 确认子应用 `render` 函数中路由 `base` 正确使用 `props.routerBase`
  - 主应用路由添加 fallback 配置（以 Vue Router 为例）：
```javascript
// 主应用路由
const router = new VueRouter({
  mode: 'history',
  routes: [
    // 其他主应用路由
    {
      path: '/vue-app/*', // 匹配子应用所有路径
      component: () => import('@/views/SubAppContainer.vue'), // 主应用挂载子应用的容器组件
    },
  ],
});
```

### 3. 样式冲突
- 原因：子应用全局样式未加前缀，或未开启 qiankun 样式隔离
- 解决：
  - 全局样式添加唯一前缀
  - 主应用注册子应用时开启 `strictStyleIsolation: true`

### 4. 子应用打包后资源路径错误
- 原因：`vue.config.js` 中 `publicPath` 配置错误（生产环境需为绝对域名）
- 解决：生产环境 `publicPath` 改为子应用线上绝对域名（如 `https://sub-vue-app.example.com/`）

## 六、总结
Vue 项目接入 qiankun 的核心是：
1. 暴露生命周期钩子（`bootstrap`/`mount`/`unmount`），兼容独立运行和嵌入模式
2. 路由 `history` 模式 + 动态 `base` 路径适配
3. 打包配置为 umd 格式 + 跨域支持
4. 样式隔离 + 主-子通信配置

改造后，子应用既可以独立开发调试，也能无缝嵌入主应用，实现微前端架构的核心目标。如果是 Vite 构建的 Vue 项目，需额外配置 Vite 的 `build.lib` 和 `server.headers`，可参考 [qiankun Vite 适配文档](https://qiankun.umijs.org/zh/guide/tutorial#vite-%E5%BA%94%E7%94%A8)。