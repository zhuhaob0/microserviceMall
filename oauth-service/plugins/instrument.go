package plugins

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
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
