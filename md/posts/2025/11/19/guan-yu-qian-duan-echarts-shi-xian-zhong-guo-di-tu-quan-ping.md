---
title: "关于前端 ECharts 实现中国地图全屏显示"
categories: [ "日常" ]
tags: [ "ECharts","中国地图","全屏显示","地图可视化","Web开发" ]
draft: false
slug: "guan-yu-qian-duan-echarts-shi-xian-zhong-guo-di-tu-quan-ping"
date: "2025-11-19 14:27:31"
url: "/guan-yu-qian-duan-echarts-shi-xian-zhong-guo-di-tu-quan-ping.html"
---

使用ECharts在前端实现中国地图全屏展示的详细步骤和关键配置：


### **核心思路**
ECharts通过内置的地图数据（或自定义地图JSON）绘制地图，实现全屏需满足两点：  
1. 地图容器占满整个浏览器可视区域（宽高设为100%）。  
2. 监听窗口大小变化，动态调整图表尺寸，确保适配屏幕。


### **实现步骤**

#### **1. 引入ECharts及地图数据**
首先需要引入ECharts库和中国地图的JSON数据（ECharts 5+版本需手动引入地图数据，不再内置）。

```html
<!-- 引入ECharts -->
<script src="https://cdn.tailwindcss.com"></script>
<script src="https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"></script>
<!-- 引入中国地图JSON数据（可从ECharts官方仓库或CDN获取） -->
<script src="https://cdn.jsdelivr.net/npm/echarts/map/json/china.json"></script>
```


#### **2. 配置全屏容器**
创建一个占满整个屏幕的容器（无滚动条，覆盖全屏）。

```html
<style>
  /* 清除默认边距，确保容器全屏 */
  * {
    margin: 0;
    padding: 0;
  }
  body, html {
    width: 100%;
    height: 100%;
    overflow: hidden; /* 隐藏滚动条 */
  }
  #mapContainer {
    width: 100%;
    height: 100%;
  }
</style>

<!-- 地图容器 -->
<div id="mapContainer"></div>
```


#### **3. 初始化ECharts地图并设置全屏配置**
通过ECharts初始化图表，配置地图参数，并监听窗口 resize 事件以自适应屏幕。

```javascript
// 初始化图表实例
const chartDom = document.getElementById('mapContainer');
const myChart = echarts.init(chartDom);

// 配置地图选项
const option = {
  // 地图类型：使用引入的中国地图数据
  series: [
    {
      type: 'map',
      map: 'china', // 对应引入的china.json的name字段
      roam: true, // 允许缩放和平移（可选，增强交互）
      label: {
        show: true, // 显示省份标签
        fontSize: 12,
        color: '#333'
      },
      itemStyle: {
        areaColor: '#f0f8ff', // 区域默认颜色
        borderColor: '#444', // 边界颜色
        borderWidth: 1
      },
      emphasis: {
        // 鼠标 hover 时的样式
        itemStyle: {
          areaColor: '#ffcc00'
        }
      }
    }
  ]
};

// 设置配置项并渲染
myChart.setOption(option);

// 监听窗口大小变化，自动调整图表尺寸
window.addEventListener('resize', () => {
  myChart.resize(); // 关键：调用resize方法适配新尺寸
});
```


### **关键说明**

1. **地图数据来源**  
   - ECharts官方已不再维护地图数据，推荐从[阿里云DataV地图选择器](https://datav.aliyun.com/portal/school/atlas/area_selector)下载最新的中国地图JSON（需注意数据完整性和边界准确性）。  
   - 下载后可本地引入，或通过CDN加载，确保`map`配置项与JSON中的`name`字段一致（通常为`china`）。

2. **全屏适配细节**  
   - 容器的`width`和`height`必须设为`100%`，且父元素（`html`、`body`）也需设置`100%`高度，否则容器无法撑开全屏。  
   - `overflow: hidden` 用于避免窗口缩放时出现滚动条，确保地图完全填充可视区域。

3. **交互增强（可选）**  
   - 添加`roam: true`允许用户缩放、平移地图，增强体验。  
   - 可通过`geo`组件替代`series.map`，支持更复杂的图层叠加（如散点、热力图等）：
     ```javascript
     option = {
       geo: {
         type: 'map',
         map: 'china',
         full: true, // 铺满容器
         roam: true
       },
       series: [] // 可添加其他系列数据（如散点）
     };
     ```

4. **常见问题**  
   - 地图显示不全：检查容器父元素是否设置`height: 100%`，或手动指定像素高度（如`height: 100vh`）。  
   - 地图无数据：确认地图JSON加载成功，且`map`配置项与JSON中的`name`一致。  
   - 窗口缩放后地图错位：确保已绑定`resize`事件并调用`myChart.resize()`。


### **扩展：添加全屏按钮**
如需手动触发全屏（类似浏览器全屏功能），可结合HTML5 Fullscreen API：

```html
<button id="fullscreenBtn">全屏显示</button>

<script>
  document.getElementById('fullscreenBtn').addEventListener('click', () => {
    const container = document.getElementById('mapContainer');
    // 进入全屏
    if (container.requestFullscreen) {
      container.requestFullscreen();
    } else if (container.webkitRequestFullscreen) { // Safari
      container.webkitRequestFullscreen();
    } else if (container.msRequestFullscreen) { // IE/Edge
      container.msRequestFullscreen();
    }
    // 全屏后调整图表尺寸
    setTimeout(() => myChart.resize(), 100);
  });
</script>
```


通过以上配置，即可实现ECharts中国地图在前端的全屏显示，且能自适应不同屏幕尺寸。根据需求可进一步优化样式、添加数据可视化效果（如省份数据标注、动态热力等）。