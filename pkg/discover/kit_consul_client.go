package discover

import (
	"final-design/pkg/common"
	"fmt"
	"log"
	"strconv"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

func New(consulHost string, consulPort string) DiscoveryClient {
	port, _ := strconv.Atoi(consulPort)
	// 通过consulHost和consulPort创建consul.Client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + consulPort
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil
	}
	client := consul.NewClient(apiClient)

	return &DiscoveryClientInstance{
		Host:   consulHost,
		Port:   port,
		config: consulConfig,
		client: client,
	}
}

func (consulClient *DiscoveryClientInstance) Register(instanceId, svcHost, healthCheckUrl, svcPort, svcName string,
	weight int, meta map[string]string, tags []string, logger *log.Logger) bool {
	port, _ := strconv.Atoi(svcPort)
	// 1. 构建服务实例元数据
	fmt.Println("Register函数: weight=", weight)
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    svcName,
		Address: svcHost,
		Port:    port,
		Meta:    meta,
		Tags:    tags,
		Weights: &api.AgentWeights{
			Passing: weight,
		},
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + svcHost + ":" + strconv.Itoa(port) + healthCheckUrl,
			Interval:                       "15s",
		},
	}

	// 2. 发送服务注册到 consul中
	err := consulClient.client.Register(serviceRegistration)
	if err != nil {
		if logger != nil {
			logger.Println("Register Service Error:", err)
		}
		return false
	}
	if logger != nil {
		logger.Println("Register Service Success!")
	}
	return true
}

func (consulClient *DiscoveryClientInstance) DeRegister(instanceId string, logger *log.Logger) bool {
	//1. 构建包含服务实例ID的元数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId,
	}

	// 2. 发送服务注销请求
	err := consulClient.client.Deregister(serviceRegistration)
	if err != nil {
		if logger != nil {
			logger.Println("DeRegister Service Error:", err)
		}
		return false
	}
	if logger != nil {
		logger.Println("DeRegister Service Success!")
	}
	return true
}

func (consulClient *DiscoveryClientInstance) DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance {
	// 该服务已监控并缓存
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}

	// 申请锁
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()
	// 再次检查是否监控
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	} else {
		// 注册监控
		go func() {
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)

			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return // 数据异常，忽略
				}
				// 没有服务实例在线
				if len(v) == 0 {
					consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
					return
				}

				var healthServices []*common.ServiceInstance
				for _, service := range v {
					if service.Checks.AggregatedStatus() == api.HealthPassing {
						healthServices = append(healthServices, newServiceInstance(service.Service))
					}
				}
				consulClient.instancesMap.Store(serviceName, healthServices)
			}

			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	}

	// 根据服务名请求服务实例列表
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
		if logger != nil {
			logger.Println("Discover Service Error:", err)
		}
		return nil
	}
	instances := make([]*common.ServiceInstance, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = newServiceInstance(entries[i].Service)
	}
	consulClient.instancesMap.Store(serviceName, instances)

	return instances
}

func newServiceInstance(service *api.AgentService) *common.ServiceInstance {
	rpcPort := service.Port - 1 // 这里如果meta没有数据的话，就默认rpcPort是servicePort减1，这个在配置文件中固定了
	fmt.Println("before===>port:", service.Port, "rpcPort:", rpcPort)
	if service.Meta != nil {
		if rpcPortString, ok := service.Meta["rpcPort"]; ok {
			rpcPort, _ = strconv.Atoi(rpcPortString)
		}
	}
	fmt.Println("after===>port:", service.Port, "rpcPort:", rpcPort)
	return &common.ServiceInstance{
		Host:     service.Address,
		Port:     service.Port,
		GrpcPort: rpcPort,
		Weight:   service.Weights.Passing,
	}
}