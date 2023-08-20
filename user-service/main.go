package main

import (
	"context"
	"final-design/pb"
	"final-design/pkg/bootstrap"
	"final-design/pkg/mysql"
	"final-design/user-service/endpoint"
	"final-design/user-service/plugins"
	"final-design/user-service/service"
	"final-design/user-service/transport"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	localconfig "final-design/user-service/config"

	register "final-design/pkg/discover"

	conf "final-design/pkg/config"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	var (
		servicePort = flag.String("service.port", bootstrap.HttpConfig.Port, "service port")
		grpcAddr    = flag.String("grpc", ":9008", "gRPC listen address")
	)
	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 5000)
	var svc service.Service = service.UserService{}

	// add logging middleware
	svc = plugins.LoggingMiddleware(localconfig.Logger)(svc)
	svc = plugins.Metrics(requestCount, requestLatency)(svc)

	userPoint := endpoint.MakeUserEndpoint(svc)
	userPoint = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(userPoint)
	userPoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "user-endpoint")(userPoint)

	createUserPoint := endpoint.MakeCreateUserEndpoint(svc)
	createUserPoint = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(createUserPoint)
	createUserPoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "create-user-endpoint")(createUserPoint)

	adminUserPoint := endpoint.MakeAdminUserEndpoint(svc)
	adminUserPoint = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(adminUserPoint)
	adminUserPoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "admin-user-endpoint")(adminUserPoint)

	createAdminUserPoint := endpoint.MakeCreateAdminUserEndpoint(svc)
	createAdminUserPoint = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(createAdminUserPoint)
	createAdminUserPoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "create-admin-user-endpoint")(createAdminUserPoint)

	// 创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)
	healthEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "health-endpoint")(healthEndpoint)

	endpts := endpoint.UserEndpoints{
		CreateUserEndpoint:      createUserPoint,
		UserEndpoint:            userPoint,
		CreateAdminUserEndpoint: createAdminUserPoint,
		AdminUserEndpoint:       adminUserPoint,

		HealthCheckEndpoint: healthEndpoint,
	}

	// 创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, localconfig.ZipkinTracer, localconfig.Logger)

	// http server
	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		mysql.InitMysql(conf.MysqlConfig.Host, conf.MysqlConfig.Port,
			conf.MysqlConfig.User, conf.MysqlConfig.Pwd, conf.MysqlConfig.Db)

		//启动前执行注册
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	// grpc server
	go func() {
		fmt.Println("grpc Server start at port:" + *grpcAddr)
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errChan <- err
			return
		}
		serverTracer := kitzipkin.GRPCServerTrace(localconfig.ZipkinTracer, kitzipkin.Name("grpc-transport"))
		tr := localconfig.ZipkinTracer
		md := metadata.MD{}
		parentSpan := tr.StartSpan("test")

		b3.InjectGRPC(&md)(parentSpan.Context())

		ctx := metadata.NewIncomingContext(context.Background(), md)
		handler := transport.NewGRPCServer(ctx, endpts, serverTracer)
		gRPCServer := grpc.NewServer()
		pb.RegisterUserServiceServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	register.Deregister()
	fmt.Println(error)
}
