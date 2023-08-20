package plugins

import (
	"final-design/sk-app/model"
	"final-design/sk-app/service"
	"time"

	"github.com/go-kit/log"
)

// loggingMiddleware Make a new type
// that contains Service interface ans logger instance
type skAppLoggingMiddleware struct {
	service.Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func SkAppLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAppLoggingMiddleware{next, logger}
	}
}

func (mw skAppLoggingMiddleware) HealthCheck() (result bool) {
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

func (mw skAppLoggingMiddleware) SecInfo(productId int) map[string]interface{} {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "SecInfo",
			"took", time.Since(begin),
		)
	}(time.Now())

	ret := mw.Service.SecInfo(productId)
	return ret
}

func (mw skAppLoggingMiddleware) SecInfoList() ([]map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "SecInfoList",
			"took", time.Since(begin),
		)
	}(time.Now())

	data, num, err := mw.Service.SecInfoList()
	return data, num, err
}

func (mw skAppLoggingMiddleware) SecKill(req *model.SecRequest) (map[string]interface{}, int, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "SecKill",
			"took", time.Since(begin),
		)
	}(time.Now())

	result, num, err := mw.Service.SecKill(req)
	return result, num, err
}
