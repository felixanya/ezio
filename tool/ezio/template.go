package main

const (
	_tplChangeLog = `## {{.Name}}

### v1.0.0
1. 上线功能xxx
`
	_tplMain = `package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/service/grpc"
	prom_wrapper "github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing_go "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/zouyx/agollo"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"{{.Name}}/internal"
	up "{{.Name}}/proto"
)

const ServiceName = "go.micro.svc.{{.Name}}"

func main() {
	// 程序启动的第一步先加载配置
	readyConfig := &agollo.AppConfig{}
	agollo.InitCustomConfig(func() (config *agollo.AppConfig, e error) {
		// 获取当前程序目录
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		e = configor.New(&configor.Config{Verbose: true}).Load(readyConfig, dir+"/app.json")
		return readyConfig, e
	})
	agollo.Start()

	// Jaeger追踪
	addr := agollo.GetStringValue("jaeger-remote-addr", "")
	log.Println("start to trace with jaeger remote addr: ", addr)
	t, io, err := NewTracer(ServiceName, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()

	// 指定Registry
	reg := consul.NewRegistry(func(op *registry.Options){
		op.Addrs = []string{
			agollo.GetStringValue("consul-addr", ""),
		}
	})

	service := grpc.NewService(
		micro.Name(ServiceName),
		micro.Version("latest"),
		// 指定服务发现地址
		micro.Registry(reg),
		micro.WrapHandler(
			// 监控
			prom_wrapper.NewHandlerWrapper(),
			// 追踪
			opentracing.NewHandlerWrapper(t),
			// 限流
			ratelimit.NewHandlerWrapper(agollo.GetIntValue("rate-limit", 0)),
		),
		micro.WrapCall(
			opentracing.NewCallWrapper(t),
		),
		micro.WrapSubscriber(
			opentracing.NewSubscriberWrapper(t),
		),
	)

	{{.Name}}Server := internal.New()

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	service.Init()

	go RegisterNewMetrics()

	err = up.RegisterUserHandler(service.Server(), {{.Name}}Server)
	if err != nil {
		fmt.Println(err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

// 输出http export
func RegisterNewMetrics() {
	addr := agollo.GetStringValue("prometheus-metrics-addr", "")
	log.Println("start to export prometheus metrics on addr: ", addr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

// 初始化opentracing，目前采用jaeger来做追踪服务
func NewTracer(serviceName, addr string) (opentracing_go.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)

	return tracer, closer, err
}
`

	_tplContributors = `# Owner
{{.Owner}}

# Author

# Reviewer
`

	_tplReadme = `# {{.Name}}

## 项目简介
1.
`
	_tplModelDefault = "package model" + `
import "time"

type Default struct {
	Id        uint64   ` + "`xorm:\"pk autoincr\"`" + `
	CreatedAt time.Time` + "`xorm:\"created\"`" + `
	UpdatedAt time.Time` + "`xorm:\"updated\"`" + `
	DeletedAt *time.Time` + "`xorm:\"deleted\"`" + `
}
`
	_tplGoMod = `module {{.Name}}

go 1.12

require (
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/coocood/freecache v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.4
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/configor v1.1.0
	github.com/micro/go-micro v1.7.1-0.20190627135301-d8e998ad85fe
	github.com/micro/go-plugins v1.1.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.6.0
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	github.com/valeamoris/ezio v0.0.0-20190712092719-dd1cfeac75e8
	github.com/zouyx/agollo v1.6.4
)
`
	_tplServer = `package internal

import (
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	myRedis "github.com/valeamoris/ezio/pkg/cache/redis"
	"github.com/valeamoris/ezio/pkg/database/mysql"
	"github.com/zouyx/agollo"
)

type Service struct {
	db    *xorm.EngineGroup
	redis *redis.Pool
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 记录错误日志
func logWrapErr(err error, format string, args ...interface{}) error {
	log.Error(errors.Wrapf(err, format, args...))
	return err
}

func New() *Service {
	var (
		mc = new(mysql.Config)
		rc = new(myRedis.Config)
	)
	// 从apollo获取数据库配置
	mysqlApollo := agollo.GetStringValue("mysql", "{}")
	fmt.Println(mysqlApollo)
	checkErr(json.Unmarshal([]byte(mysqlApollo), mc))
	redisApollo := agollo.GetStringValue("redis", "{}"")
	checkErr(json.Unmarshal([]byte(redisApollo), rc))
	return newService(mc, rc)
}

func newService(mc *mysql.Config, rc *myRedis.Config) *Service {
	return &Service{
		db:    mysql.NewMySQL(mc),
		redis: myRedis.NewPool(rc),
	}
}
`
	_tplHandler = `package internal

import (
	"context"
	"{{.Name}}/proto"
)
func (s Service) Hello(ctx context.Context, request *{{.Name}}.Request, response *{{.Name}}.Response) error {
	return nil
}
`
	_tplModel = `package model

import {{.Name}} "{{.Name}}/proto"

type Ezio struct {
	Default ` + "`xorm:\"extends\"`" + ` 
}
`
	_tplProto = `syntax = "proto3";

package {{.Name}};

message Request {}

message Response {}

service {{.Name}} {
    rpc Hello(Request) returns(Response) {}
}
`
	_tplAppJson = `
{
  "appId": "service-{{.Name}}",
  "cluster": "dev",
  "namespaceName": "application",
  "ip": "192.168.3.248:8071",
  "backupConfigPath": ""
}
`
	_tplMakefile = `.PNONY: gen-proto-go gen-proto-php build-linux

PROTO_SRC = ./proto
DST_SRC = $(PROTO_SRC)
PHP_DST_PATH = ./sdk/php

gen-proto-go:
	@echo "start to generate proto code for golang"
	protoc --proto_path=$(PROTO_SRC) --micro_out=$(DST_SRC) --go_out=$(DST_SRC) $(PROTO_SRC)/*.proto


gen-proto-php:
	@echo "start to generate proto code for php"
	protoc --proto_path=$(PROTO_SRC) --php_out=$(PHP_DST_PATH) --grpc_out=$(PHP_DST_PATH) --plugin=protoc-gen-grpc=./sdk/php/grpc_php_plugin $(PROTO_SRC)/*.proto

build-linux:
	GOOS=linux GOARCH=amd64
	go build -o bin/{{.Name}} cmd/main.go

build-docker:
	docker build -t service/{{.Name}}:latest .
`
	_tplDockerfile = `FROM golang:1.12.7 AS builder
RUN mkdir /app
WORKDIR /app
COPY . .
ENV GOPROXY https://goproxy.io
RUN ls -al
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o bin/{{.Name}} cmd/main.go

# Final image.
FROM scratch
COPY --from=builder /app/bin/{{.Name}} .
COPY --from=builder /app/app.json .
ENV MICRO_REGISTRY consul
CMD ["/{{.Name}}"]
`
)
