package plugins

import (
	"final-design/sk-admin/model"
	"final-design/sk-admin/service"
	"time"

	"github.com/go-kit/log"
	"github.com/gohouse/gorose/v2"
)

// ==============================================实现Service接口和中间件=======================================================
// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type skAdminLoggingMiddleware struct {
	service.Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func SkAdminLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
		return skAdminLoggingMiddleware{next, logger}
	}
}
func (mw skAdminLoggingMiddleware) HealthCheck() (result bool) {
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

// ==============================================实现ProductService接口和中间件=======================================================
type productLoggingMiddleware struct {
	service.ProductService
	logger log.Logger
}

func ProductLoggingMiddleware(logger log.Logger) service.ProductServiceMiddleware {
	return func(next service.ProductService) service.ProductService {
		return productLoggingMiddleware{next, logger}
	}
}

func (mw productLoggingMiddleware) CreateProduct(product *model.Product) error {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "CreateProduct",
			"product", product.ProductName,
			"took", time.Since(begin),
		)
	}(time.Now())

	err := mw.ProductService.CreateProduct(product)
	return err
}

func (mw productLoggingMiddleware) GetProductList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "GetProductList",
			"took", time.Since(begin),
		)
	}(time.Now())

	data, err := mw.ProductService.GetProductList()
	return data, err
}

func (mw productLoggingMiddleware) UpdateProduct(product *model.Product) error {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "UpdateProduct",
			"took", time.Since(begin),
		)
	}(time.Now())

	err := mw.ProductService.UpdateProduct(product)
	return err
}

func (mw productLoggingMiddleware) DeleteProduct(product *model.Product) error {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "DeleteProduct",
			"took", time.Since(begin),
		)
	}(time.Now())

	err := mw.ProductService.DeleteProduct(product)
	return err
}

// ==============================================实现ActivityService接口和中间件=======================================================
type activityLoggingMiddleware struct {
	service.ActivityService
	logger log.Logger
}

func ActivityLoggingMiddleware(logger log.Logger) service.ActivityServiceMiddleware {
	return func(next service.ActivityService) service.ActivityService {
		return activityLoggingMiddleware{next, logger}
	}
}

func (mw activityLoggingMiddleware) CreateActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "CreateActivity",
			"took", time.Since(begin),
		)
	}(time.Now())

	err := mw.ActivityService.CreateActivity(activity)
	return err
}

func (mw activityLoggingMiddleware) GetActivityList() ([]gorose.Data, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "GetActivityList",
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err := mw.ActivityService.GetActivityList()
	return ret, err
}

func (mw activityLoggingMiddleware) UpdateActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "UpdateActivity",
			"took", time.Since(begin),
		)
	}(time.Now())

	err := mw.ActivityService.UpdateActivity(activity)
	return err
}

func (mw activityLoggingMiddleware) DeleteActivity(activity *model.Activity) error {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "DeleteActivity",
			"took", time.Since(begin),
		)
	}(time.Now())

	err := mw.ActivityService.DeleteActivity(activity)
	return err
}

// ==============================================实现OrderService接口和中间件=======================================================
type orderLoggingMiddleware struct {
	service.OrderService
	logger log.Logger
}

func OrderLoggingMiddleware(logger log.Logger) service.OrderServiceMiddleware {
	return func(next service.OrderService) service.OrderService {
		return orderLoggingMiddleware{next, logger}
	}
}

func (mw orderLoggingMiddleware) GetOrderList() (map[string]interface{}, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "GetOrderList",
			"took", time.Since(begin),
		)
	}(time.Now())

	ret, err := mw.OrderService.GetOrderList()
	return ret, err
}
