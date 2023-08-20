package discover

import (
	"final-design/pkg/common"
	"log"
	"sync"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

type DiscoveryClient interface {
	/**
	服务注册接口
	@param serviceName    服务名
	@param instanceId     服务实例Id
	@param instancePort   服务实例端口
	@param healthCheckUrl 健康检查地址
	@param weight         权重
	@param meta           服务实例元数据
	*/
	Register(instanceId, svcHost, healthCheckUrl, svcPort, svcName string, weight int, meta map[string]string, tags []string, logger *log.Logger) bool

	/**
	服务注销接口
	@param instanceId 服务实例Id
	*/
	DeRegister(instanceId string, logger *log.Logger) bool

	/**
	服务实例发现接口
	@param serviceName 服务名
	*/
	DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance
}

type DiscoveryClientInstance struct {
	Host         string
	Port         int
	config       *api.Config // 连接consul的配置
	client       consul.Client
	mutex        sync.Mutex
	instancesMap sync.Map // 服务实例缓存字段
}
