package plugins

import (
	"context"
	"errors"
	"final-design/user-service/service"
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
type metricMiddleware struct {
	service.Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// Metrics 封装监控方法
func Metrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(s service.Service) service.Service {
		return metricMiddleware{
			Service:        s,
			requestCount:   requestCount,
			requestLatency: requestLatency,
		}
	}
}

func (mw metricMiddleware) Create(ctx context.Context, username, password string, userId int64, age int) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "Create"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := mw.Service.Create(ctx, username, password, userId, age)
	return err
}

func (mw metricMiddleware) Check(ctx context.Context, a, b string) (ret int64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Check"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ret, err = mw.Service.Check(ctx, a, b)
	return ret, err
}

func (mw metricMiddleware) CreateAdmin(ctx context.Context, username, password string, userId int64, age int) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateAdmin"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err := mw.Service.CreateAdmin(ctx, username, password, userId, age)
	return err
}

func (mw metricMiddleware) CheckAdmin(ctx context.Context, a, b string) (ret int64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CheckAdmin"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	ret, err = mw.Service.CheckAdmin(ctx, a, b)
	return ret, err
}

func (mw metricMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}
