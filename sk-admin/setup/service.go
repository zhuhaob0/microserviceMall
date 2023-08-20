package setup

import (
	"context"
	"final-design/sk-admin/config"
	"final-design/sk-admin/endpoint"
	"final-design/sk-admin/plugins"
	"final-design/sk-admin/service"
	"final-design/sk-admin/transport"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"final-design/pkg/bootstrap"
	"final-design/pkg/discover"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)

// 初始化Http服务
func InitSever(serviceHost, servicePort, consulHost, consulPort string) {
	log.Println("sk-admin ==> InitServer ==> servicePort is ", servicePort)
	log.Println("sk-admin ==> InitServer ==> consulPort  is ", consulPort)
	flag.Parse()

	errChan := make(chan error)
	fieldKeys := []string{"method"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "aoho",
		Subsystem: "user_user",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 1000)
	var (
		activityService service.ActivityService = service.ActivityServiceImpl{}
		productService  service.ProductService  = service.ProductServiceImpl{}
		skAdminService  service.Service         = service.SkAdminService{}
		orderService    service.OrderService    = &service.OrderServiceImpl{}
	)

	// add logging middleware
	skAdminService = plugins.SkAdminLoggingMiddleware(config.Logger)(skAdminService)
	skAdminService = plugins.SkAdminMetrics(requestCount, requestLatency)(skAdminService)

	activityService = plugins.ActivityLoggingMiddleware(config.Logger)(activityService)
	activityService = plugins.ActivityMetrics(requestCount, requestLatency)(activityService)

	productService = plugins.ProductLoggingMiddleware(config.Logger)(productService)
	productService = plugins.ProductMetrics(requestCount, requestLatency)(productService)

	orderService = plugins.OrderLoggingMiddleware(config.Logger)(orderService)
	orderService = plugins.OrderMetrics(requestCount, requestLatency)(orderService)

	// ==========================================活动endpoint========================================================
	createActivityEnd := endpoint.MakeCreateActivityEndpoint(activityService)
	createActivityEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(createActivityEnd)
	createActivityEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "create-activity")(createActivityEnd)

	GetActivityEnd := endpoint.MakeGetActivityEndpoint(activityService)
	GetActivityEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(GetActivityEnd)
	GetActivityEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-activity")(GetActivityEnd)

	updateActivityEnd := endpoint.MakeUpdateActivityEndpoint(activityService)
	updateActivityEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(updateActivityEnd)
	updateActivityEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "update-activity")(updateActivityEnd)

	deleteActivityEnd := endpoint.MakeDeleteActivityEndpoint(activityService)
	deleteActivityEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(deleteActivityEnd)
	deleteActivityEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "delete-activity")(deleteActivityEnd)

	// ==========================================商品endpoint========================================================
	createProductEnd := endpoint.MakeCreateProductEndpoint(productService)
	createProductEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(createProductEnd)
	createProductEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "create-product")(createProductEnd)

	GetProductEnd := endpoint.MakeGetProductEndpoint(productService)
	GetProductEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(GetProductEnd)
	GetProductEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-product")(GetProductEnd)

	updateProductEnd := endpoint.MakeUpdateProductEndpoint(productService)
	updateProductEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(updateProductEnd)
	updateProductEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "update-product")(updateProductEnd)

	deleteProductEnd := endpoint.MakeDeleteProductEndpoint(productService)
	deleteProductEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(deleteProductEnd)
	deleteProductEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "delete-product")(deleteProductEnd)
	// ==========================================订单endpoint========================================================
	GetOrderEnd := endpoint.MakeGetOrderEndpoint(orderService)
	GetOrderEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(GetOrderEnd)
	GetOrderEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-order")(GetOrderEnd)

	GetBuyerOrderEnd := endpoint.MakeGetBuyerOrderEndpoint(orderService)
	GetBuyerOrderEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(GetBuyerOrderEnd)
	GetBuyerOrderEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-buyer-order")(GetBuyerOrderEnd)

	// 创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(skAdminService)
	healthEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "health-endpoint")(healthEndpoint)

	endpts := endpoint.SkAdminEndpoints{
		CreateActivityEndpoint: createActivityEnd,
		GetActivityEndpoint:    GetActivityEnd,
		UpdateActivityEndpoint: updateActivityEnd,
		DeleteActivityEndpoint: deleteActivityEnd,

		CreateProductEndpoint: createProductEnd,
		GetProductEndpoint:    GetProductEnd,
		UpdateProductEndpoint: updateProductEnd,
		DeleteProductEndpoint: deleteProductEnd,

		GetOrderListEndpoint:  GetOrderEnd,
		GetBuyerOrderEndpoint: GetBuyerOrderEnd,

		HealthCheckEndpoint: healthEndpoint,
	}

	ctx := context.Background()
	// 创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.ZipkinTracer, config.Logger)

	var discoveryClient discover.DiscoveryClient = discover.New(consulHost, consulPort)
	var (
		instanceId  = bootstrap.DiscoverConfig.InstanceId
		serviceName = bootstrap.DiscoverConfig.ServiceName
		weight      = bootstrap.DiscoverConfig.Weight
	)
	//http server
	go func() {
		fmt.Println("Http Server start at port:", servicePort)
		//启动前执行注册
		if !discoveryClient.Register(instanceId, serviceHost, "/health", servicePort, serviceName, weight, nil, nil, config.Log_logger) {
			fmt.Println("discoveryClient.Register 注册失败喽~~")
			os.Exit(-1)
		}
		// fmt.Println("注册成功")
		handler := r
		errChan <- http.ListenAndServe(":"+servicePort, handler)
		fmt.Println("sk-admin 服务启动成功")
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	discoveryClient.DeRegister("skAdmin", config.Log_logger)
	fmt.Println(error)
}
