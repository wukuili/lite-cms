---
title: "apache storm demo示例"
categories: [ "hadoop","storm" ]
tags: [ "storm","apache" ]
draft: false
slug: "apache-storm-demoshili"
date: "2018-05-07 18:04:12"
url: "/apache-storm-demoshili.html"
---

从国外网站上翻译的，主要业务是创建移动电话日志分析器。
<h2>场景 - 移动呼叫日志分析器</h2>
移动电话及其持续时间将作为Apache Storm的输入提供，Storm将处理并分组相同呼叫者和接收者之间的呼叫及其呼叫总数。
<h2>创建Spout</h2>
Spout是用于数据生成的组件。基本上，spout将实现一个IRichSpout接口。“IRichSpout”界面有以下重要方法 -
<ul class="list">
 	<li><b>open</b> - 为spout提供执行环境。执行者将运行此方法来初始化spout。</li>
 	<li><b>nextTuple</b> - 通过收集器发出生成的数据。</li>
 	<li><b>close</b> - spout将要关闭时调用此方法。</li>
 	<li><b>declareOutputFields</b> - 声明元组的输出模式。</li>
 	<li><b>ack</b> - 确认处理了特定的tuple</li>
 	<li><b>fail</b> - 指定一个特定的tuple不被处理并且不被重新处理。</li>
</ul>
<h3>open</h3>
<b>open</b>方法的签名如下 -
<pre class="result notranslate">open(Map conf, TopologyContext context, SpoutOutputCollector collector)
</pre>
<ul class="list">
 	<li><b>conf</b> - 为此spout提供storm暴配置。</li>
 	<li><b>context</b> - 提供关于topology中spout位置，其任务ID，输入和输出信息的完整信息。</li>
 	<li><b>collector</b> - 使我们能够发出将由bolts处理的tuple。</li>
</ul>
<h3>nextTuple</h3>
<b>nextTuple</b>方法的签名如下 -
<pre class="result notranslate">nextTuple()
</pre>
nextTuple（）从与ack（）和fail（）方法相同的循环周期性地调用。当没有工作要做时，它必须释放对线程的控制，以便其他方法有机会被调用。所以nextTuple的第一行检查处理是否完成。如果是这样，它应该睡眠至少一毫秒，以在返回之前减少处理器上的负载。
<h3></h3>
<!--more-->
<h3>close</h3>
<b>close</b>方法的签名如下 -
<pre class="result notranslate">close()
</pre>
<h3>declareOutputFields</h3>
<b>declareOutputFields</b>方法的签名如下所示 -
<pre class="result notranslate">declareOutputFields(OutputFieldsDeclarer declarer)
</pre>
<b>declarer</b> - 它用于声明输出流ID，输出字段等。

此方法用于指定tuple的输出模式。
<h3>ack</h3>
<b>ack</b>方法的签名如下 -
<pre class="result notranslate">ack(Object msgId)
</pre>
该方法确认已经处理了特定的元组。
<h3>fail</h3>
<b>fail</b>方法的签名如下 -
<pre class="result notranslate">fail(Object msgId)
</pre>
此方法通知某个特定的元组尚未完全处理。Storm将重新处理特定的元组。
<h3>FakeCallLogReaderSpout</h3>
在我们的场景中，我们需要收集通话记录详细信息。通话记录的信息包含。
<ul class="list">
 	<li>来电号码</li>
 	<li>接收器号码</li>
 	<li>持续时间</li>
</ul>
由于我们没有实时的通话记录信息，我们会生成虚假的通话记录。假信息将使用Random类创建。完整的程序代码如下。
<h3>编码 - FakeCallLogReaderSpout.java</h3>

[code lang="java"]
package spout;
/**
 * @author: liyj
 * @Date:Created in 2018/5/8
 */
import java.util.*;
//import storm tuple packages


//import Spout interface packages

import org.apache.storm.spout.SpoutOutputCollector;
import org.apache.storm.task.TopologyContext;
import org.apache.storm.topology.IRichSpout;
import org.apache.storm.topology.OutputFieldsDeclarer;
import org.apache.storm.tuple.Fields;
import org.apache.storm.tuple.Values;

//Create a class FakeLogReaderSpout which implement IRichSpout interface to access functionalities

public class FakeCallLogReaderSpout implements IRichSpout {
    //Create instance for SpoutOutputCollector which passes tuples to bolt.
    private SpoutOutputCollector collector;
    private boolean completed = false;

    //Create instance for TopologyContext which contains topology data.
    private TopologyContext context;

    //Create instance for Random class.
    private Random randomGenerator = new Random();
    private Integer idx = 0;

    @Override
    public void open(Map conf, TopologyContext context, SpoutOutputCollector collector) {
        this.context = context;
        this.collector = collector;
    }

    @Override
    public void nextTuple() {
        if(this.idx &lt;= 1000) {
            List&lt;String&gt; mobileNumbers = new ArrayList&lt;String&gt;();
            mobileNumbers.add(&quot;1234123401&quot;);
            mobileNumbers.add(&quot;1234123402&quot;);
            mobileNumbers.add(&quot;1234123403&quot;);
            mobileNumbers.add(&quot;1234123404&quot;);

            Integer localIdx = 0;
            while(localIdx++ &lt; 100 &amp;&amp; this.idx++ &lt; 1000) {
                String fromMobileNumber = mobileNumbers.get(randomGenerator.nextInt(4));
                String toMobileNumber = mobileNumbers.get(randomGenerator.nextInt(4));

                while(fromMobileNumber == toMobileNumber) {
                    toMobileNumber = mobileNumbers.get(randomGenerator.nextInt(4));
                }

                Integer duration = randomGenerator.nextInt(60);
                this.collector.emit(new Values(fromMobileNumber, toMobileNumber, duration));
            }
        }
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        declarer.declare(new Fields(&quot;from&quot;, &quot;to&quot;, &quot;duration&quot;));
    }

    //Override all the interface methods
    @Override
    public void close() {}

    public boolean isDistributed() {
        return false;
    }

    @Override
    public void activate() {}

    @Override
    public void deactivate() {}

    @Override
    public void ack(Object msgId) {}

    @Override
    public void fail(Object msgId) {}

    @Override
    public Map&lt;String, Object&gt; getComponentConfiguration() {
        return null;
    }
}

[/code]

&nbsp;
<h2>创建Bolt</h2>
Bolt是一个将tuple作为输入，处理tuple并生成新的tuple作为输出的组件。Bolts将实现<b>IRichBolt</b>接口。在这个程序中，使用两个螺栓类<b>CallLogCreatorBolt</b>和<b>CallLogCounterBolt</b>来执行操作。

IRichBolt接口有以下方法 -
<ul class="list">
 	<li><b>prepare</b> - 为bolt提供执行的环境。执行者将运行此方法来初始化spout。</li>
 	<li><b>execute</b> - 处理输入的单个tuple。</li>
 	<li><b>cleanup</b> - 当bolt即将关闭时调用。</li>
 	<li><b>declareOutputFields</b> - 声明tuple的输出模式。</li>
</ul>
<h3>准备</h3>
<b>准备</b>方法的签名如下 -
<pre class="result notranslate">prepare(Map conf, TopologyContext context, OutputCollector collector)
</pre>
<ul class="list">
 	<li><b>conf</b> - 为此bolt提供storm配置。</li>
 	<li><b>context</b> - 提供有关topology中bolt位置，其任务ID，输入和输出信息等的完整信息。</li>
 	<li><b>collector</b> - 使我们能够发出处理过的tuple。</li>
</ul>
<h3>execute</h3>
<b>execute</b>方法的签名如下 -
<pre class="result notranslate">execute(Tuple tuple)
</pre>
这里的<b>tuple</b>是要处理的输入tuple。

所述<b>execute</b>方法一次处理单tuple。tuple数据可以通过Tuple类的getValue方法访问。没有必要立即处理输入tuple。多tuple可以作为单个输出tuple进行处理和输出。处理过的tuple可以通过使用OutputCollector类发出。
<h3>cleanup</h3>
<b>cleanup</b>方法的签名如下 -
<pre class="result notranslate">cleanup()
</pre>
<h3>declareOutputFields</h3>
<b>declareOutputFields</b>方法的签名如下所示 -
<pre class="result notranslate">declareOutputFields(OutputFieldsDeclarer declarer)
</pre>
这里参数<b>declarer</b>用于声明输出流ID，输出字段等。

此方法用于指定tuple的输出模式
<h2>通话记录创建者bolt</h2>
通话记录创建器bolt接收通话记录tuple。通话记录tuple具有主叫号码，接收者号码和通话时长。通过组合主叫方号码和接收方号码，此bolt简单地创建一个新值。新值的格式为“来电号码 - 接收方号码”，并将其命名为新字段“call”。完整的代码如下。
<h3>编码 - CallLogCreatorBolt.java</h3>

[code lang="java"]
package bolt;

/**
 * @author: liyj
 * @Date:Created in 2018/5/8
 */
//import util packages
import org.apache.storm.task.OutputCollector;
import org.apache.storm.task.TopologyContext;
import org.apache.storm.topology.IRichBolt;
import org.apache.storm.topology.OutputFieldsDeclarer;
import org.apache.storm.tuple.Fields;
import org.apache.storm.tuple.Tuple;
import org.apache.storm.tuple.Values;

import java.util.HashMap;
import java.util.Map;



//Create a class CallLogCreatorBolt which implement IRichBolt interface
public class CallLogCreatorBolt implements IRichBolt {
    //Create instance for OutputCollector which collects and emits tuples to produce output
    private OutputCollector collector;

    @Override
    public void prepare(Map conf, TopologyContext context, OutputCollector collector) {
        this.collector = collector;
    }

    @Override
    public void execute(Tuple tuple) {
        String from = tuple.getString(0);
        String to = tuple.getString(1);
        Integer duration = tuple.getInteger(2);
        collector.emit(new Values(from + &quot; - &quot; + to, duration));
    }

    @Override
    public void cleanup() {}

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        declarer.declare(new Fields(&quot;call&quot;, &quot;duration&quot;));
    }

    @Override
    public Map&lt;String, Object&gt; getComponentConfiguration() {
        return null;
    }
}

[/code]

<h2>通话记录计数器bolt</h2>
呼叫记录计数器bolt接收呼叫及其持续时间作为tuple。这个bolt在prepare方法中初始化一个字典（Map）对象。在<b>execute</b>方法中，它检查tuple并在tuple中为每个新的“call”值在字典对象中创建一个新条目，并在字典对象中设置值1。对于字典中已有的条目，它只是递增其值。简单地说，这个bolt将call和它的计数保存在字典对象中。我们可以将它保存到数据源中，而不是将call和它的计数保存在字典中。完整的程序代码如下所示 -
<h3>编码 - CallLogCounterBolt.java</h3>

[code lang="java"]
package bolt;

/**
 * @author: liyj
 * @Date:Created in 2018/5/8
 */
import org.apache.storm.task.OutputCollector;
import org.apache.storm.task.TopologyContext;
import org.apache.storm.topology.IRichBolt;
import org.apache.storm.topology.OutputFieldsDeclarer;
import org.apache.storm.tuple.Fields;
import org.apache.storm.tuple.Tuple;

import java.util.HashMap;
import java.util.Map;



public class CallLogCounterBolt implements IRichBolt {
    Map&lt;String, Integer&gt; counterMap;
    private OutputCollector collector;

    @Override
    public void prepare(Map conf, TopologyContext context, OutputCollector collector) {
        this.counterMap = new HashMap&lt;String, Integer&gt;();
        this.collector = collector;
    }

    @Override
    public void execute(Tuple tuple) {
        String call = tuple.getString(0);
        Integer duration = tuple.getInteger(1);

        if(!counterMap.containsKey(call)){
            counterMap.put(call, 1);
        }else{
            Integer c = counterMap.get(call) + 1;
            counterMap.put(call, c);
        }

        collector.ack(tuple);
    }

    @Override
    public void cleanup() {
        for(Map.Entry&lt;String, Integer&gt; entry:counterMap.entrySet()){
            System.out.println(entry.getKey()+&quot; : &quot; + entry.getValue());
        }
    }

    @Override
    public void declareOutputFields(OutputFieldsDeclarer declarer) {
        declarer.declare(new Fields(&quot;call&quot;));
    }

    @Override
    public Map&lt;String, Object&gt; getComponentConfiguration() {
        return null;
    }

}
[/code]

<h2>创建topology</h2>
Storm topology基本上是一个Thrift结构。TopologyBuilder类提供了简单和轻松的方法来创建复杂的topology。TopologyBuilder类具有设置spout<b>（setSpout）</b>和设置bolt<b>（setBolt）的方法</b>。最后，TopologyBuilder创建Topology来创建topology。使用下面的代码片段来创建一个拓扑 -

[code lang="java"]
TopologyBuilder builder = new TopologyBuilder();

builder.setSpout(&quot;call-log-reader-spout&quot;, new FakeCallLogReaderSpout());

builder.setBolt(&quot;call-log-creator-bolt&quot;, new CallLogCreatorBolt())
   .shuffleGrouping(&quot;call-log-reader-spout&quot;);

builder.setBolt(&quot;call-log-counter-bolt&quot;, new CallLogCounterBolt())
   .fieldsGrouping(&quot;call-log-creator-bolt&quot;, new Fields(&quot;call&quot;));
[/code]

<b>shuffleGrouping</b>和<b>fieldsGrouping</b>方法有助于设置spout和bolt的流分组。
<h2>本地集群</h2>
出于开发目的，我们可以使用“LocalCluster”对象创建本地集群，然后使用“LocalCluster”类的“submitTopology”方法提交topology。“submitTopology”的一个参数是“Config”类的一个实例。在提交topology之前，“Config”类用于设置配置选项。该配置选项将在运行时与集群配置合并，并通过prepare方法发送到所有任务（spout和bolt）。将topology提交到集群后，我们将等待10秒钟，以便集群计算提交的topology，然后使用“LocalCluster”的“close”方法关闭群集。完整的程序代码如下所示 -
<h3>编码 - LogAnalyserStorm.java</h3>

[code lang="java"]
package topology;

import bolt.CallLogCounterBolt;
import bolt.CallLogCreatorBolt;
import org.apache.storm.Config;
import org.apache.storm.LocalCluster;
import org.apache.storm.topology.TopologyBuilder;
import org.apache.storm.tuple.Fields;
import spout.FakeCallLogReaderSpout;

/**
 * @author: liyj
 * @Date:Created in 2018/5/8
 */


//Create main class LogAnalyserStorm submit topology.
public class LogAnalyserStorm {
    public static void main(String[] args) throws Exception{
        //Create Config instance for cluster configuration
        Config config = new Config();
        config.setDebug(true);

        //
        TopologyBuilder builder = new TopologyBuilder();
        builder.setSpout(&quot;call-log-reader-spout&quot;, new FakeCallLogReaderSpout());

        builder.setBolt(&quot;call-log-creator-bolt&quot;, new CallLogCreatorBolt())
                .shuffleGrouping(&quot;call-log-reader-spout&quot;);

        builder.setBolt(&quot;call-log-counter-bolt&quot;, new CallLogCounterBolt())
                .fieldsGrouping(&quot;call-log-creator-bolt&quot;, new Fields(&quot;call&quot;));

        LocalCluster cluster = new LocalCluster();
        cluster.submitTopology(&quot;LogAnalyserStorm&quot;, config, builder.createTopology());
        Thread.sleep(10000);

        //Stop the topology

        cluster.shutdown();
    }
}
[/code]

<h2>构建和运行应用程序(我已经将程序编写为maven项目，大家可以将项目clone到本地在自己的ide中查看)</h2>
完整的应用程序有四个Java代码。他们是 -
<ul class="list">
 	<li>FakeCallLogReaderSpout.java</li>
 	<li>CallLogCreaterBolt.java</li>
 	<li>CallLogCounterBolt.java</li>
 	<li>LogAnalyerStorm.java</li>
</ul>
应用程序可以使用以下命令构建 -
<pre class="result notranslate">javac -cp “/path/to/storm/apache-storm-0.9.5/lib/*” *.java
</pre>
应用程序可以使用以下命令运行 -
<pre class="result notranslate">java -cp “/path/to/storm/apache-storm-0.9.5/lib/*”:. LogAnalyserStorm</pre>
我在github上写的demo,直接clone到本地，在ide中打开即可。

https://github.com/wukuili/storm_test.git
<h3>产量</h3>
一旦应用程序启动，它将输出有关集群启动过程，spout和bolt处理的完整详细信息，最后还会输出集群关闭过程。在“CallLogCounterBolt”中，我们打印了call及其计数详细信息。这些信息将如下显示在控制台上 -
<pre class="result notranslate">1234123402 - 1234123401 : 78
1234123402 - 1234123404 : 88
1234123402 - 1234123403 : 105
1234123401 - 1234123404 : 74
1234123401 - 1234123403 : 81
1234123401 - 1234123402 : 81
1234123403 - 1234123404 : 86
1234123404 - 1234123401 : 63
1234123404 - 1234123402 : 82
1234123403 - 1234123402 : 83
1234123404 - 1234123403 : 86
1234123403 - 1234123401 : 93
</pre>
<h2>非JVM语言</h2>
Storm风格的topologies结构通过Thrift接口实现，这使得用任何语言提交topologies变得非常容易。Storm支持Ruby，Python和许多其他语言。我们来看看python绑定。
<h3>Python绑定</h3>
Python是一种通用的解释型，交互式，面向对象和高级编程语言。Storm支持Python来实现其topology。Python支持emitting, anchoring, acking和记录操作。

如你所知，bolt可以用任何语言来定义。以另一种语言编写的bolt作为子流程执行，Storm通过标准输入/标准输出与JSON消息通信。首先拿一个支持python绑定的示例bolt WordCount。

[code lang="python"]
public static class WordCount implements IRichBolt {
   public WordSplit() {
      super(&quot;python&quot;, &quot;splitword.py&quot;);
   }
	
   public void declareOutputFields(OutputFieldsDeclarer declarer) {
      declarer.declare(new Fields(&quot;word&quot;));
   }
}

[/code]

这里的<b>WordCount</b>类实现了<b>IRichBolt</b>接口，并使用python实现指定的超级方法参数“splitword.py”运行。现在创建一个名为“splitword.py”的python实现。

[code lang="python"]
import storm
   class WordCountBolt(storm.BasicBolt):
      def process(self, tup):
         words = tup.values[0].split(&quot; &quot;)
         for word in words:
         storm.emit([word])
WordCountBolt().run()
[/code]

这是用于计算给定句子中的单词数量的Python示例实现。同样，您也可以使用其他支持语言进行绑定。