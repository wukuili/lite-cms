---
title: "rabbitMq常用创建消息应用的maven demo项目（一）---路由routing"
categories: [ "hadoop","rabbitMq" ]
tags: [ "rabbitmq","routing" ]
draft: false
slug: "rabbitmqchangyongchuangjianxiaoxiyingyongdemaven-demoxiangmuyi-luyourouting"
date: "2018-05-14 21:42:59"
url: "/rabbitmqchangyongchuangjianxiaoxiyingyongdemaven-demoxiangmuyi-luyourouting.html"
---

rabbitmq官网上提供了6个demo，分别从是hello world、工作队列、发布/订阅、路由、主题、rpc这六个demo。
基本上看完这6哥demo之后，对rabbitmq应该就有了清晰的认识，并且可以达到基本数量应用的程度。

下面我挑选最常用的路由和主题这两个demo，为大家翻译下。个人加谷歌翻译，有不合适的地方，欢迎大家批评指正。

Routing---路由

在之前的教程中，我们构建了一个简单的日志系统 我们能够将日志消息广播给许多接收者。

<span class="goog-text-highlight">在本教程中，我们将在他的基础上添加一个功能 - 只订阅一部分消息。</span>例如，我们只将严重错误的消息导入日志文件（以节省磁盘空间），同时仍然能够在控制台上打印所有日志消息。

<!--more-->
<h2>绑定</h2>
在前面的例子中，我们已经创建了绑定。您可能会回想一下代码：
<pre>[code lang="java"]channel.queueBind（queueName，EXCHANGE_NAME，“”）;[/code]</pre>
binding是exchange和queue之间的桥梁。可以简单地理解为：queue对来自该exchange的消息感兴趣。

bindings可以采用额外的<span class="code ">routingKey</span>参数。为了避免与<span class="code ">basic_publish</span>参数混淆，我们将其称为 <span class="code ">binding key</span>。下面就是我们如何使用一个key创建一个bindings：
<pre>[code lang="java"]channel.queueBind（queueName，EXCHANGE_NAME，“black”）;[/code]</pre>
binding key的含义取决于exchange的类型。但此key对于exchange为fanout的类型无效。（因为fanout类型的exchange是将消息发给全部queue）
<h2>direct exchange</h2>
我们之前教程的日志记录系统将所有消息广播给所有消费者。我们希望将其扩展一个功能：可以根据消息的严重性进行过滤。例如，我们可能需要一个将日志消息中仅仅是严重错误的写入磁盘，而不会在warn或info级别的日志消息中浪费磁盘空间。

我们正在使用一个<span class="code ">fanout exchange</span>，这没有给我们很大的灵活性 - 它只能进行盲目的广播。

我们将使用<span class="code ">direct exchange</span>。<span class="code ">direct exchange</span>背后的路由算法很简单 - 消息进入队列，其 <span class="code ">绑定密钥</span>与消息的<span class="code ">路由密钥</span>完全匹配。

为了说明这一点，请参考下面的图：
<div class="diagram"> <img class="alignnone size-medium wp-image-137" src="https://lequ7.com/wp-content/uploads/2018/05/direct-exchange-300x126.png" alt="" width="300" height="126" /></div>
在这个图中，我们可以看到有两个队列绑定的<span class="code ">direct exchange </span><span class="code ">X. </span>第一个队列用bindingKey:orange绑定，第二个队列有两个绑定，一个bindingKey为<span class="code ">black</span>，另一个为<span class="code ">green</span>。

在这种图中，使用routing将<span class="code ">orange</span>发布到exchange的消息 将被路由到队列<span class="code ">Q1</span>。带有<span class="code ">black</span> 或<span class="code ">green</span>路由键的消息将进入<span class="code ">Q2</span>。所有其他消息将被丢弃。
<h2>多个绑定multiply bindings</h2>
<div class="diagram"><img src="http://www.rabbitmq.com/img/tutorials/direct-exchange-multiple.png" height="170" /></div>
使用相同的bindingKey绑定多个queue是完全合法的。在我们的例子中，我们可以使用绑定键<span class="code ">black</span>添加<span class="code ">X</span>和<span class="code ">Q1</span>之间的绑定。在这种情况下，<span class="code ">直接</span>交换就像<span class="code ">fanout类型</span>一样，将消息广播到所有匹配的队列。带有路由键<span class="code ">black的</span>消息将传送到 <span class="code ">Q1</span>和<span class="code ">Q2</span>。
<h2>发布日志Emitting Logs</h2>
我们将把这个模型用于我们的日志系统。这次我们不用fanout的exchange,而是将消息发送到<span class="code ">direck exchange</span>。我们将日志严重级别作为<span class="code ">key</span>。这样接收程序将能够选择想要接收的严重程度。我们先关注发布日志。

与往常一样，我们需要先创建一个exchange：
<pre class="sourcecode java hljs">channel.exchangeDeclare（EXCHANGE_NAME，<span class="hljs-string">“direct”</span>）;
</pre>
我们准备发送一条消息：
<pre class="sourcecode java hljs">channel.basicPublish（EXCHANGE_NAME，severity，<span class="hljs-keyword">null，</span>message.getBytes（））;
</pre>
为了简化，我们将假设“严重级别”可以是'info'，'warning'，'error'之一。
<h2>订阅Subscribing</h2>
接收消息的方式与上一个教程中的一样，只有一个例外 - 我们将为每个我们感兴趣的严重级别创建一个新绑定。
<pre class="sourcecode java hljs">String queueName = channel.queueDeclare（）.getQueue（）;

<span class="hljs-keyword">for</span>（String severity：argv）{
  channel.queueBind（queueName，EXCHANGE_NAME，severity）;
}
</pre>
<h2>把它放在一起</h2>
<div class="diagram"><img src="http://www.rabbitmq.com/img/tutorials/python-four.png" height="170" /></div>
<span class="code ">EmitLogDirect.java</span>类的代码：
<pre>[code lang="java"]
import com.rabbitmq.client.*;

import java.io.IOException;

public class EmitLogDirect {

    private static final String EXCHANGE_NAME = &quot;direct_logs&quot;;

    public static void main(String[] argv)
                  throws java.io.IOException {

        ConnectionFactory factory = new ConnectionFactory();
        factory.setHost(&quot;localhost&quot;);
        Connection connection = factory.newConnection();
        Channel channel = connection.createChannel();

        channel.exchangeDeclare(EXCHANGE_NAME, &quot;direct&quot;);

        String severity = getSeverity(argv);
        String message = getMessage(argv);

        channel.basicPublish(EXCHANGE_NAME, severity, null, message.getBytes());
        System.out.println(&quot; [x] Sent '&quot; + severity + &quot;':'&quot; + message + &quot;'&quot;);

        channel.close();
        connection.close();
    }
    //..
}
[/code]</pre>
ReceiveLogsDirect.java的代码:
<pre>[code lang="java"]
import com.rabbitmq.client.*;

import java.io.IOException;

public class ReceiveLogsDirect {

  private static final String EXCHANGE_NAME = &quot;direct_logs&quot;;

  public static void main(String[] argv) throws Exception {
    ConnectionFactory factory = new ConnectionFactory();
    factory.setHost(&quot;localhost&quot;);
    Connection connection = factory.newConnection();
    Channel channel = connection.createChannel();

    channel.exchangeDeclare(EXCHANGE_NAME, &quot;direct&quot;);
    String queueName = channel.queueDeclare().getQueue();

    if (argv.length &amp;amp;lt; 1){
      System.err.println(&quot;Usage: ReceiveLogsDirect [info] [warning] [error]&quot;);
      System.exit(1);
    }

    for(String severity : argv){
      channel.queueBind(queueName, EXCHANGE_NAME, severity);
    }
    System.out.println(&quot; [*] Waiting for messages. To exit press CTRL+C&quot;);

    Consumer consumer = new DefaultConsumer(channel) {
      @Override
      public void handleDelivery(String consumerTag, Envelope envelope,
                                 AMQP.BasicProperties properties, byte[] body) throws IOException {
        String message = new String(body, &quot;UTF-8&quot;);
        System.out.println(&quot; [x] Received '&quot; + envelope.getRoutingKey() + &quot;':'&quot; + message + &quot;'&quot;);
      }
    };
    channel.basicConsume(queueName, true, consumer);
  }
}
[/code]</pre>
接下来官网里面是用java命令分别运行这两个class文件。

我们如果是在ide中的话，直接运行即可（如果直接run的话，需要配置下参数）

<img class="alignnone wp-image-136" src="https://lequ7.com/wp-content/uploads/2018/05/rabbitmq-300x190.png" alt="" width="376" height="238" />