package common

// ServiceInstance 服务实例，拥有以下属性
type ServiceInstance struct {
	Host      string
	Port      int
	Weight    int
	CurWeight int
	GrpcPort  int
}
