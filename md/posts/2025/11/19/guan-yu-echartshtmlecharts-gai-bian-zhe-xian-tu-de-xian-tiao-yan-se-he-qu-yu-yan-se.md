---
title: "关于 ECharts 在 HTML 中修改折线图的线条颜色和区域颜色"
categories: [ "日常" ]
tags: [ "ECharts","折线图","线条颜色","区域颜色","HTML" ]
draft: false
slug: "guan-yu-echartshtmlecharts-gai-bian-zhe-xian-tu-de-xian-tiao-yan-se-he-qu-yu-yan-se"
date: "2025-11-19 14:29:11"
url: "/guan-yu-echartshtmlecharts-gai-bian-zhe-xian-tu-de-xian-tiao-yan-se-he-qu-yu-yan-se.html"
---

在ECharts中，折线图（`line`类型）的线条颜色和区域填充颜色可通过`itemStyle`和`areaStyle`配置项灵活设置，以下是具体实现方法和示例：


### **核心配置说明**
折线图的颜色设置主要涉及两个配置：
- **线条颜色**：通过 `series.itemStyle.color` 或 `series.lineStyle.color` 定义（推荐用`lineStyle`更直观）。
- **区域颜色**：通过 `series.areaStyle` 配置，可设置填充色、透明度等，让折线与X轴之间的区域显示颜色。


### **完整示例代码**
以下是一个包含线条颜色和区域颜色设置的折线图完整HTML示例：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <title>ECharts折线图颜色设置</title>
  <!-- 引入ECharts -->
  <script src="https://cdn.jsdelivr.net/npm/echarts@5.4.3/dist/echarts.min.js"></script>
  <style>
    #lineChart {
      width: 800px;
      height: 500px;
      margin: 20px auto;
    }
  </style>
</head>
<body>
  <!-- 图表容器 -->
  <div id="lineChart"></div>

  <script>
    // 初始化图表
    const chartDom = document.getElementById('lineChart');
    const myChart = echarts.init(chartDom);

    // 配置项
    const option = {
      xAxis: {
        type: 'category',
        data: ['1月', '2月', '3月', '4月', '5月', '6月']
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: '销量',
          type: 'line', // 折线图类型
          data: [120, 190, 130, 150, 220, 280],
          smooth: true, // 线条平滑（可选）
          
          // 1. 设置线条颜色
          lineStyle: {
            color: '#36CFC9', // 线条主色（支持十六进制、RGB、颜色名）
            width: 3, // 线条粗细
            type: 'solid' // 线条类型：solid/dashed/dotted
          },
          
          // 2. 设置区域填充颜色（折线与X轴之间的区域）
          areaStyle: {
            color: {
              // 渐变填充（可选，比单一颜色更美观）
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: '#36CFC9' }, // 区域顶部颜色（与线条同色）
                { offset: 1, color: 'rgba(54, 207, 201, 0.1)' } // 区域底部颜色（透明渐变）
              ]
            }
          },
          
          // 3. 数据点颜色（可选，默认继承线条颜色）
          itemStyle: {
            color: '#1890FF', // 数据点颜色
            borderColor: '#fff', // 数据点边框色
            borderWidth: 2 // 数据点边框粗细
          }
        }
      ]
    };

    // 渲染图表
    myChart.setOption(option);

    // 窗口大小变化时自适应
    window.addEventListener('resize', () => {
      myChart.resize();
    });
  </script>
</body>
</html>
```


### **关键配置详解**
1. **线条颜色（`lineStyle`）**  
   - `color`：支持多种格式，例如：
     - 十六进制：`#36CFC9`
     - RGB：`rgb(54, 207, 201)`
     - RGBA（带透明度）：`rgba(54, 207, 201, 0.8)`
     - 颜色名：`'teal'`
   - `width`：线条粗细（像素值，默认2）。
   - `type`：线条样式（`solid`实线、`dashed`虚线、`dotted`点线）。

2. **区域颜色（`areaStyle`）**  
   - 单一颜色：直接设置`color: 'rgba(54, 207, 201, 0.3)'`。
   - 渐变颜色（推荐）：通过`colorStops`定义渐变起点和终点，`offset`范围0~1（0为顶部，1为底部），实现自然过渡效果。
   - 注意：需确保`series`中没有禁用区域填充（默认不设置`areaStyle`则无填充，设置后才会显示区域颜色）。

3. **多系列折线图颜色设置**  
   若图表包含多条折线，可在`series`数组中为每个系列单独配置`lineStyle`和`areaStyle`，示例：
   ```javascript
   series: [
     {
       name: '产品A',
       type: 'line',
       data: [120, 190, 130],
       lineStyle: { color: '#36CFC9' },
       areaStyle: { color: 'rgba(54, 207, 201, 0.2)' }
     },
     {
       name: '产品B',
       type: 'line',
       data: [90, 150, 200],
       lineStyle: { color: '#FF4D4F' },
       areaStyle: { color: 'rgba(255, 77, 79, 0.2)' }
     }
   ]
   ```


### **常见问题**
- **区域颜色不显示**：检查是否遗漏`areaStyle`配置，或`areaStyle.color`是否设置了全透明（如`rgba(0,0,0,0)`）。
- **线条颜色与数据点不一致**：`itemStyle.color`优先级高于`lineStyle.color`，如需保持一致可省略`itemStyle.color`，默认继承线条颜色。
- **渐变效果失效**：确保渐变配置的`type`（`linear`线性/`radial`径向）和坐标（`x,y,x2,y2`）正确，线性渐变中`x2,y2`决定方向（如`y2:1`表示从上到下渐变）。


通过以上配置，可灵活定制折线图的线条和区域颜色，适配不同的视觉需求。如需更复杂的效果（如动态变色、条件颜色），可结合ECharts的`visualMap`组件或回调函数实现。