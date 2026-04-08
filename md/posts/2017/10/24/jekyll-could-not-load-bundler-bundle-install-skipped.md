---
title: "jekyll Could not load Bundler. Bundle install skipped."
categories: [ "jekyll" ]
tags: [ "jekyll" ]
draft: false
slug: "jekyll-could-not-load-bundler-bundle-install-skipped"
date: "2017-10-24 14:56:49"
url: "/jekyll-could-not-load-bundler-bundle-install-skipped.html"
---

<blockquote>安装Jekyll的时候出现错误，提示Could not load Bundler. Bundle install skipped. 然后使用Jekyll serve的时候提示 没有运行所需的bunddle,从stackoverflow上找到的答案</blockquote>
安装Jekyll的时候出现错误，提示<code class="highlighter-rouge">Could not load Bundler. Bundle install skipped.</code> 然后使用Jekyll serve的时候提示 没有运行所需的bunddle,从stackoverflow上找到的答案
<figure class="highlight">
<pre><code class="language-ruby" data-lang="ruby"><span class="n">gem</span> <span class="n">install</span> <span class="n">jekyll</span> <span class="n">bundler</span> 
<span class="n">jekyll</span> <span class="n">new</span> <span class="n">my</span><span class="o">-</span><span class="n">awesome</span><span class="o">-</span><span class="n">site</span> 
<span class="n">cd</span> <span class="n">my</span><span class="o">-</span><span class="n">awesome</span><span class="o">-</span><span class="n">site</span> 
<span class="n">bundle</span> <span class="n">install</span> 
<span class="n">bundle</span> <span class="nb">exec</span> <span class="n">jekyll</span> <span class="n">serve</span>
<span class="c1"># =&gt; Now browse to http://localhost:4000</span></code></pre>
</figure>