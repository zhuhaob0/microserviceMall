package plugins

import (
	"context"
	"errors"
	"final-design/sk-app/model"
	"final-design/sk-app/service"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/juju/ratelimit"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("rate limit exceed")

// NewTokenBucketLimitterWithJUju 使用juju/ratelimit创建限流中间件
func NewTokenBucketLimitterWithJUju(bkt *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if bkt.TakeAvailable(1) == 0 { // 桶中没有可用令牌返回0，不会阻塞
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

// NewTokenBucketLimitterWithBuildIn 使用x/time/rate创建限流中间件
func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

// metricMiddleware 定义监控中间件，嵌入Service
// 新增监控指标：requestCount和requestLatency
type skAppMetricMiddleware struct {
	service.Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// Metrics 封装监控方法
func SkAppMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(s service.Service) service.Service {
		return skAppMetricMiddleware{
			Service:        s,
			requestCount:   requestCount,
			requestLatency: requestLatency,
		}
	}
}

func (mw skAppMetricMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}

func (mw skAppMetricMiddleware) SecInfo(productId int) map[string]interface{} {
	defer func(begin time.Time) {
		lvs := []string{"method", "SecInfo"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ret := mw.Service.SecInfo(productId)
	return ret
}

func (mw skAppMetricMiddleware) SecInfoList() ([]map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SecInfoList"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, num, err := mw.Service.SecInfoList()
	return data, num, err
}

func (mw skAppMetricMiddleware) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SecKill"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result, num, err := mw.Service.SecKill(req)
	return result, num, err
}
