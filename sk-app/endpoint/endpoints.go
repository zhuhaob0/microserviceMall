package endpoint

import (
	"context"
	"final-design/sk-app/model"
	"final-design/sk-app/service"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

// CalculationEndpoint define endpoint
type SkAppEndpoints struct {
	SecKillEndpoint        endpoint.Endpoint
	HealthCheckEndpoint    endpoint.Endpoint
	GetSecInfoEndpoint     endpoint.Endpoint
	GetSecInfoListEndpoint endpoint.Endpoint
	TestEndpoint           endpoint.Endpoint
}

func (ue SkAppEndpoints) HealthCheck() bool {
	return false
}

type SecInfoRequest struct {
	ProductId int `json:"product_id"`
}

type Response struct {
	Result map[string]interface{} `json:"result"`
	Error  string                 `json:"error"`
	Code   int                    `json:"code"`
}

type SecInfoListResponse struct {
	Result []map[string]interface{} `json:"result"`
	Num    int                      `json:"num"`
	Error  string                   `json:"error"`
}

// make endpoint
func MakeSecInfoEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SecInfoRequest)
		ret := svc.SecInfo(req.ProductId)
		if ret == nil {
			err = fmt.Errorf("productId=[%d]的商品不存在", req.ProductId)
			return Response{Result: ret, Error: err.Error()}, nil
		}
		return Response{Result: ret, Error: ""}, nil
	}
}

func MakeSecInfoListEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ret, num, err := svc.SecInfoList()
		if err != nil {
			return SecInfoListResponse{ret, num, err.Error()}, nil
		}
		return SecInfoListResponse{ret, num, ""}, nil
	}
}

// make endpoint
func MakeSecKillEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.SecRequest)
		ret, code, calError := svc.SecKill(&req)

		if calError != nil {
			return Response{Result: ret, Code: code, Error: calError.Error()}, nil
		}
		return Response{Result: ret, Code: code, Error: ""}, nil
	}
}

func MakeTestEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return Response{Result: nil, Code: 1, Error: ""}, nil
	}
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

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
