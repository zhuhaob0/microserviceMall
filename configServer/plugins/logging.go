package plugins

import (
	"final-design/configServer/service"
	"time"

	"github.com/go-kit/log"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
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

func (mw loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthChcek",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	result = mw.Service.HealthCheck()
	return
}

func (mw loggingMiddleware) DownloadFile(name string) (content []byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "DownloadFile",
			"result", len(content) != 0,
			"took", time.Since(begin),
		)
	}(time.Now())
	content, err = mw.Service.DownloadFile(name)
	return
}
