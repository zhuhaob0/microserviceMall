package loadbalance

import (
	"errors"
	"final-design/pkg/common"
	"math/rand"
)

var (
	ErrServiceInstanceNotExist = errors.New("service instances are not exist")
)

// 负载均衡器
type LoadBalance interface {
	SelectService(services []*common.ServiceInstance) (*common.ServiceInstance, error)
}

type RandomLoadBalance struct {
}

// 随机负载均衡
func (loadbalance *RandomLoadBalance) SelectService(services []*common.ServiceInstance) (*common.ServiceInstance, error) {
	if len(services) == 0 { // services 等于nil时，len(services)也等于0
		return nil, ErrServiceInstanceNotExist
	}

	return services[rand.Intn(len(services))], nil
}

type WeightRoundRobinLoadBalance struct {
}

// 权重平滑负载均衡
func (loadbalance *WeightRoundRobinLoadBalance) SelectService(services []*common.ServiceInstance) (best *common.ServiceInstance, err error) {
	if len(services) == 0 {
		return nil, ErrServiceInstanceNotExist
	}
	total := 0
	for i := 0; i < len(services); i++ {
		w := services[i]
		if w == nil {
			continue
		}
		w.CurWeight += w.Weight
		total += w.Weight
		if best == nil || w.CurWeight > best.CurWeight {
			best = w
		}
	}

	if best == nil {
		return nil, nil
	}
	best.CurWeight -= total
	return best, nil
}
