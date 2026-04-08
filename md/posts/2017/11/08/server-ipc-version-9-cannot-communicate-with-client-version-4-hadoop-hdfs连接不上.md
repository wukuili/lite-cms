---
title: "Server IPC version 9 cannot communicate with client version 4  hadoop hdfs连接不上"
categories: [ "hadoop" ]
tags: [ "hdfs","hadoop" ]
draft: false
slug: "server-ipc-version-9-cannot-communicate-with-client-version-4-hadoop-hdfs连接不上"
date: "2017-11-08 09:00:55"
url: "/server-ipc-version-9-cannot-communicate-with-client-version-4-hadoop-hdfs连接不上.html"
---

<blockquote>commons-httpclient-3.1.jarcommons-io-2.4.jarcommons-lang-2.6.jarcommons-logging-1.1.3.jarcommons-net-3.1.jarguava-11.0.2.jarhadoop-common-2.6.2.jarhadoop-auth-2.6.2.jarslf4j-api-1.7.5.jarhadoop-hdfs-2.6.2.jarcommons-cli-1.2.jarprotobuf-java-2.5.0.jarhtrace-core-3.0.4.jar</blockquote>
在pom.xml中添加这些

commons-httpclient-3.1.jar
commons-io-2.4.jar
commons-lang-2.6.jar
commons-logging-1.1.3.jar
commons-net-3.1.jar
guava-11.0.2.jar
hadoop-common-2.6.2.jar
hadoop-auth-2.6.2.jar
slf4j-api-1.7.5.jar
hadoop-hdfs-2.6.2.jar
commons-cli-1.2.jar
protobuf-java-2.5.0.jar
htrace-core-3.0.4.jar

&nbsp;

以下为示例代码：

&nbsp;
<code>
/**
 * @Autohor: liyj
 * @Description: <p/>
 * @Date:Created in 2017/11/7
 * @Modified by :
 */

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FileStatus;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;

public class HdfsFileReader {

    private static final String NAME_NODE = "hdfs://tj02:8020";//nameNomeHost = localhost if you use hadoop in local mode

    public static void main(String[] args) throws URISyntaxException, IOException {
//        String fileInHdfs = "/user/hive/warehouse/t001011003";
        Configuration configuration = new Configuration();
//        configuration.set("fs.hdfs.impl",
//                org.apache.hadoop.hdfs.DistributedFileSystem.class.getName()
//        );
//        configuration.set("fs.file.impl",
//                org.apache.hadoop.fs.LocalFileSystem.class.getName()
//        );
        FileSystem fs = FileSystem.get(URI.create(NAME_NODE), configuration);
//
//        fs.createNewFile(new Path("/user/hive/warehouse/t001011003/0000sadasd"));
//                String fileContent = IOUtils.toString(fs.open(new Path(fileInHdfs)), "UTF-8");
//        System.out.println("File content - " + fileContent);
//        copyFile2Hdfs();
        Path listf = new Path("/user/hive/warehouse/t001011003");
        FileStatus stats[] = fs.listStatus(listf);

        for (int i = 0; i < stats.length; ++i)

        {

            System.out.println(stats[i].getPath().toString());

        }

        fs.close();
    }

    public static void copyFile2Hdfs() throws IOException {
        Configuration conf = new Configuration();

        FileSystem hdfs = FileSystem.get(conf);


        //本地文件

//        Path src =new Path("D:\\HebutWinOS");

        //HDFS为止

        Path dst = new Path("/");


//        hdfs.copyFromLocalFile(src, dst);

        System.out.println("Upload to" + conf.get("fs.default.name"));


        FileStatus files[] = hdfs.listStatus(dst);

        for (FileStatus file : files) {

            System.out.println(file.getPath());

        }
    }

}

</code>