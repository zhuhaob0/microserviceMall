package config

import (
	"fmt"
	"log"
	"os"

	"final-design/pkg/bootstrap"
	conf "final-design/pkg/config"

	kitlog "github.com/go-kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
)

const (
	kConfigType = "CONGFIG_TYPE"
)

var ZipkinTracer *zipkin.Tracer
var Logger kitlog.Logger
var Log_logger *log.Logger

func init() {
	fmt.Println("调用sk-admin/config.go的init函数")
	Logger = kitlog.NewLogfmtLogger(os.Stderr)
	Logger = kitlog.With(Logger, "ts", kitlog.DefaultTimestampUTC)
	Logger = kitlog.With(Logger, "caller", kitlog.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	if err := conf.LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config", err)
	}
	if err := conf.Sub("mysql", &conf.MysqlConfig); err != nil {
		Logger.Log("Fail to parse mysql", err)
	}
	if err := conf.Sub("trace", &conf.TraceConfig); err != nil {
		Logger.Log("Fail to parse trace", err)
	}

	zipkinUrl := "http://" + conf.TraceConfig.Host + ":" + conf.TraceConfig.Port + conf.TraceConfig.Url
	Logger.Log("zipkin url", zipkinUrl)
	initTracer(zipkinUrl)
}

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
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
