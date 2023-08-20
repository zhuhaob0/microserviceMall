package endpoint

import (
	"context"
	"errors"
	"final-design/sk-admin/model"
	"final-design/sk-admin/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/gohouse/gorose/v2"
)

// CalculateEndpoint define endpoint
type SkAdminEndpoints struct {
	CreateActivityEndpoint endpoint.Endpoint
	GetActivityEndpoint    endpoint.Endpoint
	UpdateActivityEndpoint endpoint.Endpoint
	DeleteActivityEndpoint endpoint.Endpoint

	CreateProductEndpoint endpoint.Endpoint
	GetProductEndpoint    endpoint.Endpoint
	UpdateProductEndpoint endpoint.Endpoint
	DeleteProductEndpoint endpoint.Endpoint

	GetOrderListEndpoint  endpoint.Endpoint
	GetBuyerOrderEndpoint endpoint.Endpoint

	HealthCheckEndpoint endpoint.Endpoint
}

func (s SkAdminEndpoints) HealthCheck() bool {
	return false
}

var (
	ErrInvalidRequestType = errors.New("invalid username, password")
)

// UserRequest define request struct
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse define response struct
type UserResponse struct{}

type GetResponse struct {
	Result []gorose.Data `json:"result"`
	Error  error         `json:"error"`
}

// 普通的response，增删改都用它
type Response struct {
	Error error `json:"error"`
}

type OrderResponse struct {
	Result map[string]interface{} `json:"result"`
	Error  error                  `json:"error"`
}

type BuyerOrderRequest struct {
	Buyer string `json:"buyer"`
}

type BuyerOrderResponse struct {
	Result []gorose.Data `json:"result"`
	Error  error         `json:"error"`
}

// ========================================================活动Endpoint===============================================

// 创建获取所有活动列表的endpoint
func MakeGetActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		activityList, calError := svc.GetActivityList()
		if calError != nil {
			return GetResponse{Result: nil, Error: calError}, nil
		}
		return GetResponse{Result: activityList, Error: calError}, nil
	}
}

// 创建新增活动的endpoint
func MakeCreateActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Activity)
		calError := svc.CreateActivity(&req)
		return Response{Error: calError}, nil
	}
}

// 创建修改活动时的endpoint
func MakeUpdateActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Activity)
		calError := svc.UpdateActivity(&req)
		return Response{Error: calError}, nil
	}
}

// 创建删除活动时的endpoint
func MakeDeleteActivityEndpoint(svc service.ActivityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Activity)
		calError := svc.DeleteActivity(&req)
		return Response{Error: calError}, nil
	}
}

// ========================================================商品Endpoint===============================================

// 创建获取所有商品列表的endpoint
func MakeGetProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getProductList, calError := svc.GetProductList()
		if calError != nil {
			return GetResponse{Result: nil, Error: calError}, nil
		}
		return GetResponse{Result: getProductList, Error: calError}, nil
	}
}

// 创建新增商品的endpoint
func MakeCreateProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Product)
		calError := svc.CreateProduct(&req)
		return Response{Error: calError}, nil
	}
}

// 创建更新商品的endpoint
func MakeUpdateProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Product)
		calError := svc.UpdateProduct(&req)
		return Response{Error: calError}, nil
	}
}

// 创建删除商品的endpoint
func MakeDeleteProductEndpoint(svc service.ProductService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.Product)
		calError := svc.DeleteProduct(&req)
		return Response{Error: calError}, nil
	}
}

// ========================================================订单Endpoint===============================================

func MakeGetOrderEndpoint(svc service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		orderList, orderErr := svc.GetOrderList()
		if orderErr != nil {
			return OrderResponse{Result: nil, Error: orderErr}, nil
		} else {
			return OrderResponse{Result: orderList, Error: orderErr}, nil
		}
	}
}

func MakeGetBuyerOrderEndpoint(svc service.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(BuyerOrderRequest)
		buyerOrder, err := svc.GetBuyerOrder(req.Buyer)
		if err != nil {
			return BuyerOrderResponse{Result: nil, Error: err}, nil
		} else {
			return BuyerOrderResponse{Result: buyerOrder, Error: err}, nil
		}
	}
}

// ========================================================健康检查Endpoint===============================================

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

type GetListRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{Status: status}, nil
	}
}
