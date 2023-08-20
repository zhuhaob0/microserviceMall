package plugins

import (
	"context"
	"final-design/user-service/service"
	"time"

	"github.com/go-kit/log"
)

// loggingMiddleware Make a new type
// that contains Service interface ans logger instance
type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) Create(ctx context.Context, username, password string, userId int64, age int) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Create",
			"result", err == nil,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.Create(ctx, username, password, userId, age)
	return
}

func (mw loggingMiddleware) Check(ctx context.Context, a, b string) (ret int64, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Check",
			"username", a,
			"pwd", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.Check(ctx, a, b)
	return ret, err
}

func (mw loggingMiddleware) CreateAdmin(ctx context.Context, username, password string, userId int64, age int) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "CreateAdmin",
			"result", err == nil,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.CreateAdmin(ctx, username, password, userId, age)
	return
}

func (mw loggingMiddleware) CheckAdmin(ctx context.Context, a, b string) (ret int64, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "CheckAdmin",
			"username", a,
			"pwd", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err = mw.Service.CheckAdmin(ctx, a, b)
	return ret, err
}

func (mw loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthCheck",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())

	result = mw.Service.HealthCheck()
	return
}
