package setup

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"final-design/sk-app/config"
	"final-design/sk-app/endpoint"
	"final-design/sk-app/plugins"
	"final-design/sk-app/service"

	// kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	// stdprometheus "github.com/prometheus/client_golang/prometheus"
	localconfig "final-design/pkg/config"

	"final-design/sk-app/transport"

	register "final-design/pkg/discover"

	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"golang.org/x/time/rate"
)

// 初始化http服务
func InitServer(host string, servicePort string) {
	log.Println("sk-app service port is ", servicePort)
	flag.Parse()

	errChan := make(chan error)

	rateBucket := rate.NewLimiter(rate.Every(time.Second), 5000)
	var skAppService service.Service = service.SkAppService{}
	skAppService = plugins.SkAppLoggingMiddleware(config.Logger)(skAppService)

	healthCheckEnd := endpoint.MakeHealthCheckEndpoint(skAppService)
	healthCheckEnd = plugins.NewTokenBucketLimitterWithBuildIn(rateBucket)(healthCheckEnd)
	healthCheckEnd = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "health-check")(healthCheckEnd)

	GetSecInfoEnd := endpoint.MakeSecInfoEndpoint(skAppService)
	GetSecInfoEnd = plugins.NewTokenBucketLimitterWithBuildIn(rateBucket)(GetSecInfoEnd)
	GetSecInfoEnd = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "sec-info")(GetSecInfoEnd)

	GetSecInfoListEnd := endpoint.MakeSecInfoListEndpoint(skAppService)
	GetSecInfoListEnd = plugins.NewTokenBucketLimitterWithBuildIn(rateBucket)(GetSecInfoListEnd)
	GetSecInfoListEnd = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "sec-info-list")(GetSecInfoListEnd)

	// 秒杀接口单独限流
	secRatebucket := rate.NewLimiter(rate.Every(time.Microsecond*100), 5000)

	SecKillEnd := endpoint.MakeSecKillEndpoint(skAppService)
	SecKillEnd = plugins.NewTokenBucketLimitterWithBuildIn(secRatebucket)(SecKillEnd)
	SecKillEnd = plugins.AuthToken()(SecKillEnd)
	// SecKillEnd = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "sec-kill")(SecKillEnd)

	testEnd := endpoint.MakeTestEndpoint(skAppService)
	testEnd = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "test")(testEnd)

	endpts := endpoint.SkAppEndpoints{
		SecKillEndpoint:        SecKillEnd,
		HealthCheckEndpoint:    healthCheckEnd,
		GetSecInfoEndpoint:     GetSecInfoEnd,
		GetSecInfoListEndpoint: GetSecInfoListEnd,
		TestEndpoint:           testEnd,
	}
	ctx := context.Background()
	// 创建http handler
	r := transport.MakeHttpHandler(ctx, endpts, localconfig.ZipkinTracer, localconfig.Logger)

	// http server
	go func() {
		fmt.Println("Http Server start at port:" + servicePort)
		//启动前执行注册
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	<-errChan
	// 服务退出取消注册
	register.Deregister()
}
