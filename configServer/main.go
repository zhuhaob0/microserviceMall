package main

import (
	"context"
	"final-design/configServer/config"
	"final-design/configServer/endpoint"
	"final-design/configServer/plugins"
	"final-design/configServer/service"
	"final-design/configServer/transport"
	"final-design/pkg/discover"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 获取命令行参数
	var (
		servicePort = flag.String("service.port", "10085", "service port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		consulPort  = flag.String("consul.port", "8500", "consul port")
		consulHost  = flag.String("consul.host", "127.0.0.1", "consul host")
		serviceName = flag.String("service.name", "config-service", "service name")
		configName  = "configService"
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)
	var discoveryClient discover.DiscoveryClient = discover.New(*consulHost, *consulPort)

	var svc service.Service = service.ConfigService{}

	// add logging middleware
	svc = plugins.LoggingMiddleware(config.KitLogger)(svc)

	// 创建文件下载Endpoint
	downloadEndpoint := endpoint.MakeDownloadFileEndpoint(svc)
	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	//把DownloadEndpoint 和HealthCheckEndpoint 封装至ConfigEndpoint
	endpts := endpoint.ConfigEndpoints{
		DownloadFileEndpoint: downloadEndpoint,
		HealthCheckEndpoint:  healthEndpoint,
	}

	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)

	instanceId := configName

	//http server
	go func() {

		config.Logger.Println("Config Server start at port:" + *servicePort)
		//启动前执行注册
		if !discoveryClient.Register(instanceId, *serviceHost, "/health", string(*servicePort), *serviceName, 10, nil, nil, config.Logger) {
			config.Logger.Printf("configService for service %s failed.", *serviceName)
			// 注册失败，服务启动失败
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	discoveryClient.DeRegister(instanceId, config.Logger)
	config.Logger.Println(error)
}
