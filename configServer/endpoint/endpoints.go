package endpoint

import (
	"context"
	"final-design/configServer/service"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

// ConfigEndpoint define endpoint
type ConfigEndpoints struct {
	DownloadFileEndpoint endpoint.Endpoint
	HealthCheckEndpoint  endpoint.Endpoint
}

type DownloadRequest struct {
	Filename string `json:"file_name"`
}
type DownloadResponse struct {
	Content []byte `json:"content"`
}

func MakeDownloadFileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DownloadRequest)
		data, err := svc.DownloadFile(req.Filename)
		if err != nil {
			fmt.Println("svc.DownloadFile 出错:", err)
			return DownloadResponse{}, err
		}
		return DownloadResponse{Content: data}, nil
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
		return HealthResponse{status}, nil
	}
}
