#ezio

##### 服务链路追踪
jaeger服务端初始化
```bash
$ docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.13
```

##### 应用程序监控
prometheus初始化
```
tar xvfz prometheus-*.tar.gz
cd prometheus-*

//添加目标 默认http协议  metrics
scrape_configs:
  - job_name:       'example-random'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:8080', 'localhost:8081']
        labels:
          group: 'production'

      - targets: ['localhost:8082']
        labels:
          group: 'canary'
          
//启动prometheus 
./prometheus --config.file=prometheus.yml
```

mac grafana初始化
```
brew update
brew install grafana

brew services start grafana
```

##### tool
这个项目是基于bilibili的kratos的tool，并且改造为适合项目的版本

##### 项目目录
通常一个services会被分成几个部分

    --cmd 项目执行程序入口
    --migrate 数据库同步入口
    --model 模型定义文件
    --proto 原型文件及生成的文件目录
    --sdk 生成的php文件及golang client目录文件
    app.json.sample 从apollo获取配置的基础配置，该配置默认和运行程序在同目录下
    Dockerfile