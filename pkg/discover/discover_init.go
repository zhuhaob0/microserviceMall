package discover

import (
	"errors"
	"final-design/pkg/bootstrap"
	"final-design/pkg/common"
	"final-design/pkg/loadbalance"
	"fmt"
	"log"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
)

var (
	ConsulService       DiscoveryClient
	LoadBalance         loadbalance.LoadBalance
	Logger              *log.Logger
	ErrNoInstanceExited = errors.New("no available client")
)

func init() {
	// 实例化一个Consul客户端，此处实例化了原生态实现版本
	ConsulService = New(bootstrap.DiscoverConfig.Host, bootstrap.DiscoverConfig.Port)
	LoadBalance = new(loadbalance.RandomLoadBalance)
	Logger = log.New(os.Stderr, "", log.LstdFlags)
}

func CheckHealth(writer http.ResponseWriter, reader *http.Request) {
	Logger.Println("Health Check!")
	_, err := fmt.Fprintln(writer, "Server is ok!")
	if err != nil {
		Logger.Println(err)
	}
}

func DiscoverService(serviceName string) (*common.ServiceInstance, error) {
	instances := ConsulService.DiscoverServices(serviceName, Logger)

	if len(instances) < 1 {
		Logger.Printf("no avaliable client for %s.", serviceName)
		return nil, ErrNoInstanceExited
	}
	return LoadBalance.SelectService(instances)
}

func Register() {
	// 实例失败，停止服务
	if ConsulService == nil {
		panic(0)
	}
	// 判空instanceId 通过go.uuid 获取一个服务实例ID
	instanceId := bootstrap.DiscoverConfig.InstanceId
	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + uuid.NewV4().String()
	}

	if !ConsulService.Register(instanceId, bootstrap.HttpConfig.Host, "/health", bootstrap.HttpConfig.Port,
		bootstrap.DiscoverConfig.ServiceName, bootstrap.DiscoverConfig.Weight,
		map[string]string{
			"rpcPort": bootstrap.RpcConfig.Port,
		}, nil, Logger) {
		Logger.Printf("register service %s failed.", bootstrap.DiscoverConfig.ServiceName)
		// 注册失败，服务启动失败
		panic(0)
	}
	Logger.Printf(bootstrap.DiscoverConfig.ServiceName+"-service for service %s success.", bootstrap.DiscoverConfig.ServiceName)
}

func Deregister() {
	// 实例失败，停止服务
	if ConsulService == nil {
		panic(0)
	}
	// 判空 instanceId 通过go.uuid获取一个服务实例Id
	instanceId := bootstrap.DiscoverConfig.InstanceId
	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.InstanceId
	}

	if !ConsulService.DeRegister(instanceId, Logger) {
		Logger.Printf("deregister for service %s failed.", bootstrap.DiscoverConfig.ServiceName)
		panic(0)
	}
}
