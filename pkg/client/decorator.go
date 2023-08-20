package client

import (
	"context"
	"errors"
	"final-design/pkg/bootstrap"
	"final-design/pkg/discover"
	"final-design/pkg/loadbalance"
	"fmt"
	"log"
	"strconv"
	"time"

	conf "final-design/pkg/config"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"google.golang.org/grpc"
)

var (
	ErrRPCService = errors.New("no rpc service")
)

var defaultLoadBalance loadbalance.LoadBalance = &loadbalance.RandomLoadBalance{}

type ClientManager interface {
	DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer,
		ctx context.Context, inputVal interface{}, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName     string
	logger          *log.Logger
	discoveryClient discover.DiscoveryClient
	loadBalance     loadbalance.LoadBalance
	after           []InvokeAfterFunc
	before          []InvokeBeforeFunc
}

type InvokeAfterFunc func() (err error)
type InvokeBeforeFunc func() (err error)

func (manager *DefaultClientManager) DecoratorInvoke(path string, hystrixName string,
	tracer opentracing.Tracer, ctx context.Context, inputVal interface{}, outVal interface{}) (err error) {

	// 1. 回调函数
	for _, fn := range manager.before {
		if err = fn(); err != nil {
			return err
		}
	}

	// 2. 使用Hystrix的Do方法构造对应的断路器保护
	err = hystrix.Do(hystrixName, func() error {
		// 3. 服务发现
		instances := manager.discoveryClient.DiscoverServices(manager.serviceName, manager.logger)
		// 4. 负载均衡
		if instance, err := manager.loadBalance.SelectService(instances); err == nil {
			fmt.Println("instance.GrpcPort=", instance.GrpcPort)
			if instance.GrpcPort > 0 {
				// 5. 获得RPC端口并且发送RPC请求
				if conn, err := grpc.Dial(instance.Host+":"+strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(),
					grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(genTracer(tracer), otgrpc.LogPayloads())),
					grpc.WithTimeout(time.Second*1)); err == nil {
					fmt.Println("path=", path)
					if err = conn.Invoke(ctx, path, inputVal, outVal); err != nil {
						fmt.Println("conn.Invoke error:", err)
						return err
					}
				} else {
					fmt.Println("grpc.Dial出错,err:", err)
					return err
				}

			} else {
				return ErrRPCService
			}

		} else {
			fmt.Println("manager.loadBalance.SelectService出错,err:", err)
			return err
		}

		return nil
	}, func(e error) error {
		fmt.Println("hystrix的fallback函数触发,err:", err)
		return e
	})

	if err != nil {
		return err
	} else {
		for _, fn := range manager.after {
			if err = fn(); err != nil {
				return err
			}
		}
		return nil
	}
}

func genTracer(tracer opentracing.Tracer) opentracing.Tracer {
	if tracer != nil {
		return tracer
	}
	zipkinUrl := "http://" + conf.TraceConfig.Host + ":" + conf.TraceConfig.Port + conf.TraceConfig.Url
	zipkinRecorder := bootstrap.HttpConfig.Host + ":" + bootstrap.HttpConfig.Port
	collector, err := zipkin.NewHTTPCollector(zipkinUrl)
	if err != nil {
		log.Fatalf("zipkin.NewHTTPCollector err:%v", err)
	}
	recorder := zipkin.NewRecorder(collector, false, zipkinRecorder, bootstrap.DiscoverConfig.ServiceName)

	res, err := zipkin.NewTracer(recorder, zipkin.ClientServerSameSpan(true))
	if err != nil {
		log.Fatalf("zipkin.NewTracer err:%v", err)
	}
	return res
}
