package plugins

import (
	"context"
	"errors"
	"final-design/sk-admin/model"
	"final-design/sk-admin/service"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/gohouse/gorose/v2"
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
type skAdminMetricMiddleware struct {
	service.Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

type activityMetricMiddleware struct {
	service.ActivityService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

type productMetricMiddleware struct {
	service.ProductService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

type orderMetricMiddleware struct {
	service.OrderService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// Metrics 封装监控方法
func SkAdminMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(s service.Service) service.Service {
		return skAdminMetricMiddleware{
			Service:        s,
			requestCount:   requestCount,
			requestLatency: requestLatency,
		}
	}
}

// Metrics 封装监控方法
func ProductMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productMetricMiddleware{
			ProductService: next,
			requestCount:   requestCount,
			requestLatency: requestLatency,
		}
	}
}

// Metrics 封装监控方法
func ActivityMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityMetricMiddleware{
			ActivityService: next,
			requestCount:    requestCount,
			requestLatency:  requestLatency,
		}
	}
}

func OrderMetrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.OrderServiceMiddleware {
	return func(next service.OrderService) service.OrderService {
		return orderMetricMiddleware{
			OrderService:   next,
			requestCount:   requestCount,
			requestLatency: requestLatency,
		}
	}
}

func (mw skAdminMetricMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}

// =========================================商品======================================================

func (mw productMetricMiddleware) CreateProduct(product *model.Product) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateProduct"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	error := mw.ProductService.CreateProduct(product)
	return error
}

func (mw productMetricMiddleware) GetProductList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetProductList"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	data, err := mw.ProductService.GetProductList()
	return data, err
}

func (mw productMetricMiddleware) UpdateProduct(product *model.Product) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateProduct"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	error := mw.ProductService.UpdateProduct(product)
	return error
}

func (mw productMetricMiddleware) DeleteProduct(product *model.Product) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteProduct"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	error := mw.ProductService.DeleteProduct(product)
	return error
}

// =========================================活动======================================================

func (mw activityMetricMiddleware) CreateActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateActivity"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	error := mw.ActivityService.CreateActivity(activity)
	return error
}

func (mw activityMetricMiddleware) GetActivityList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetActivityList"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result, err := mw.ActivityService.GetActivityList()
	return result, err
}

func (mw activityMetricMiddleware) UpdateActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateActivity"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	error := mw.ActivityService.UpdateActivity(activity)
	return error
}

func (mw activityMetricMiddleware) DeleteActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteActivity"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	error := mw.ActivityService.DeleteActivity(activity)
	return error
}

// =========================================订单======================================================

func (mw orderMetricMiddleware) GetOrderList() (map[string]interface{}, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetOrderList"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result, err := mw.OrderService.GetOrderList()
	return result, err
}
