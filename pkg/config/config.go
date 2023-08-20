package config

import (
	"bytes"
	"encoding/base64"
	"errors"
	"final-design/pkg/bootstrap"
	"final-design/pkg/discover"
	"final-design/pkg/loadbalance"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	_log "log"

	"github.com/go-kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
)

const (
	kConfigType = "CONFIG_TYPE"
)

var (
	ZipkinTracer *zipkin.Tracer
	Logger       log.Logger
	logger       *_log.Logger
	// 权重平滑负载均衡
	LoadBalance loadbalance.LoadBalance = new(loadbalance.WeightRoundRobinLoadBalance)
)

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
}

func init() {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	if err := LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config", err)
	}

	//if err := Sub("mysql", &MysqlConfig); err != nil {
	//	Logger.Log("Fail to parse mysql", err)
	//}

	if err := Sub("trace", &TraceConfig); err != nil {
		Logger.Log("Fail to parse trace", err)
		return // add by myself
	}
	zipkinUrl := "http://" + TraceConfig.Host + ":" + TraceConfig.Port + TraceConfig.Url
	Logger.Log("zipkin url", zipkinUrl)
	initTracer(zipkinUrl)
}

func initTracer(zipkinURL string) {
	var (
		err           error
		useNoopTracer = zipkinURL == ""
		reporter      = zipkinhttp.NewReporter(zipkinURL)
	)
	// defer reporter.Close()
	zEP, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, bootstrap.HttpConfig.Port)
	ZipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer))

	if err != nil {
		Logger.Log("err", err)
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinURL)
	}
}

func LoadRemoteConfig() (err error) {
	// 声明consul实例
	discoverClientInstance := discover.New(bootstrap.DiscoverConfig.Host, bootstrap.DiscoverConfig.Port)
	// 服务实例发现
	serviceInstances := discoverClientInstance.DiscoverServices(bootstrap.ConfigServerConfig.Id, logger)
	// 负载均衡算法寻找最合适的实例
	serviceInstance, err := LoadBalance.SelectService(serviceInstances)
	if err != nil {
		Logger.Log("LoadBalance.SelectService Happend An Error", err)
		return
	}
	configServer := "http://" + serviceInstance.Host + ":" + strconv.Itoa(serviceInstance.Port)
	confAddr := fmt.Sprintf("%v/%v/%v-%v.%v", configServer, bootstrap.ConfigServerConfig.Label,
		bootstrap.DiscoverConfig.ServiceName, bootstrap.ConfigServerConfig.Profile, viper.Get(kConfigType))

	fmt.Println("confAddr=", confAddr)
	resp, err := http.Get(confAddr)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取响应的内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// 将响应的内容解析为viper配置
	viper.SetConfigType(viper.GetString(kConfigType)) // 默认yaml
	err = viper.ReadConfig(bytes.NewBuffer(body))
	if err != nil {
		return
	}
	// 先读出content的具体内容，再将其解析为viper配置
	decodeStr, _ := base64.StdEncoding.DecodeString(viper.Get("content").(string))
	content := []byte(decodeStr)
	err = viper.ReadConfig(bytes.NewBuffer(content))
	if err != nil {
		return
	}
	Logger.Log("Load config from: ", confAddr)
	return
}

func Sub(key string, value interface{}) error {
	Logger.Log("配置文件的前缀：", key)
	sub := viper.Sub(key)
	if sub == nil {
		fmt.Println("./pkg/config/config.go 的Sub函数出错")
		return errors.New("sub is nil")
	}
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
