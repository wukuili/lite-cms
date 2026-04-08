---
title: "react-native flatlist上拉加载onEndReached不能正确触发的解决办法"
categories: [ "react-native" ]
tags: [ "react native","flatlist" ]
draft: false
slug: "react-native-flatlistshanglajiazaionendreachedbunengzhengquechufadejiejuebanfa"
date: "2018-08-27 17:36:09"
url: "/react-native-flatlistshanglajiazaionendreachedbunengzhengquechufadejiejuebanfa.html"
---

<!-- wp:paragraph {"backgroundColor":"pale-cyan-blue"} -->
<p class="has-background has-pale-cyan-blue-background-color"><strong>问题</strong><strong></strong></p>
<!-- /wp:paragraph -->

<!-- wp:paragraph -->
<p>在写flatlist复用组件时，调用的时候如果父组件是不定高的组件，会造成组件无法显示</p>
<!-- /wp:paragraph -->

<!-- wp:paragraph -->
<p>如果父组件样式{flex：1}，则会出现下拉方法频繁触发或不正常触发的问题（我这里出现的问题是在列表第6个项目在底部时，缓慢上拉会多次触发flatlist的onEndReached监听）</p>
<!-- /wp:paragraph -->

<!-- wp:paragraph {"backgroundColor":"vivid-green-cyan"} -->
<p class="has-background has-vivid-green-cyan-background-color">原因</p>
<!-- /wp:paragraph -->

<!-- wp:paragraph -->
<p>推测是因为{flex：1}不适合做动态高度组件的父组件样式，会错误的判断高度导致onEndReached多次不正常触发。</p>
<!-- /wp:paragraph -->

<!-- wp:paragraph {"backgroundColor":"vivid-green-cyan"} -->
<p class="has-background has-vivid-green-cyan-background-color">解决</p>
<!-- /wp:paragraph -->

<!-- wp:paragraph -->
<p>可以把列表上方所需的组件做成header属性传入组件当做flatlist的头部组件，这样就可以直接调用封装好的组件。</p>
<!-- /wp:paragraph -->

<!-- wp:paragraph -->
<p>也可以把父元素的样式设成{height: '100%'}，这样就可以正确的触发onEndReached监听。</p>
<!-- /wp:paragraph -->