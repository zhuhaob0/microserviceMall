package bootstrap

var (
	HttpConfig         HttpConf
	DiscoverConfig     DiscoverConf
	ConfigServerConfig ConfigServerConf
	RpcConfig          RpcConf
)

// Http 配置
type HttpConf struct {
	Host string
	Port string
}

// RPC配置
type RpcConf struct {
	Port string
}

// 服务注册与发现配置
type DiscoverConf struct {
	Host        string
	Port        string
	ServiceName string
	Weight      int
	InstanceId  string
}

// 配置中心
type ConfigServerConf struct {
	Id      string
	Profile string
	Label   string
}
