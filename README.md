<h2># XHYL grpc日志收集服务</h2>
XHYL技术栈：golang 1.13、elasticsearch7.4.2、 kibana7.4.2、docker部署、consul

<h3>简介</h3>
<div>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;xhyl日志收集服务是基于golang编写，部署一个grpc的日志收集服务，客户端引入xhyl client包到，在自己的demo中可使用logc进行调用Info、Debug等相关函数将日志发送到server端，server端将日志存入elasticsearch，使用kibana可视化日志数据</div>

<h4>环境安装（linux）：</h4>

1.安装docker

2.docker安装elasticsearch(docker镜像仓库：https://hub.docker.com/_/elasticsearch)

    拉去镜像：docker pull elasticsearch:7.4.2
    运行容器：docker run -d --name elasticsearch -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:7.4.2
    进入elasticsearch配置跨域：
            docker exec -it elasticsearch /bin/bash
            vi config/elasticsearch.yml
            添加配置：   http.cors.enabled: true
                    http.cors.allow-origin: "*"
    重启es容器：docker restart elasticsearch

3.docker安装kibana（docker镜像仓库：https://hub.docker.com/_/kibana）

    拉去镜像：docker pull kibana:7.4.2
    运行容器：docker run --name kibana -e ELASTICSEARCH_HOSTS=http://192.168.10.20:9200  -p 5601:5601 -d kibana:7.4.2

注意：kibana和es版本要对应，ELASTICSEARCH_HOSTS为es访问地址

4.访问kibana（浏览器访问）：http://localhost:5601


<h4>服务端部署步骤(linux部署，windows下可直接运行调试):</h4>

1.下载服务端源码:

    git clone https://github.com/xx132917/xhyl-logger.git

2.打包服务端:
    
    set GOARCH=amd64
    set GOOS=linux
    go build -o main ./

3.部署:

    将打包好的可执行文件COPY到一个服务器中，需COPY文件：config.json、Dockerfile、seelog.xml、main（打包好的可执行文件）
    ---
        config.json配置文件：
            {
              "Port": 5021,
              "EsHost": "http://192.168.10.20:9200",
            }
    ---

    chmod 777 main
    ./main

    如果你向将该服务模块注册至consul可在config.json中添加配置进行启动：
    ---
        config.json配置文件：
            {
              "Port": 5021,
              "EsHost": "http://192.168.10.20:9200",
              "ConsulAddress": "192.168.10.20:8500", //consul服务ip:port
              "ConsulRegisterId": "logserver",
              "ConsulRegisterName": "logserver",
              "ConsulRegisterPort": 5022,
              "ConsulRegisterTags": "v0001",
              "ConsulRegisterAddress": "192.168.10.20", //本日志服务部署的服务器ip
              "ConsulCheckPort": 5022,
              "ConsulCheckTimeout": "3s",
              "ConsulCheckInterval": "5s",
              "DeregisterCriticalServiceAfter": "20s"
            }
    ---

4.docker部署

    （1）cd到config.json、Dockerfile、seelog.xml、main所在目录
    （2）制作docker镜像：
                       docker build -t logconnection:v1 .
                       docker images
         注意：基础镜像为centos:latest : docker pull centos:latest (如果有特殊需求，请自选择基础镜像)
    （3）运行容器：
                docker run -d --restart=unless-stopped --name=logconnection -p 5021:5021 -p 5022:5022 -v $pwd:/home/  logconnection:v1
         注意：-v 是你挂载config.json配置文件的目录

    （4）docker ps查看容器启动成功


<h4>客户端使用:</h4>

     import (
     	"logconnection/proto/client" 
     	"time"
     )

     func main() {
            logc.SetLogcAddress("192.168.10.33:5021")
            logc.SetLogcNodeName("logcclient")
            
            //测试使用
            ogc.Error("Error 6666666666")
            
            ime.Sleep(time.Second)
     }
