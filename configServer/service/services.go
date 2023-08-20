package service

import (
	"fmt"
	"os"
)

// Service Define a service interface
type Service interface {
	// HealthCheck check service health status
	HealthCheck() bool
	DownloadFile(name string) ([]byte, error)
}

// ArithmeticService implement Service interface
type ConfigService struct {
}

func (s ConfigService) DownloadFile(name string) ([]byte, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		fmt.Printf("读取文件[%s]失败\n", name)
		return []byte{}, err
	}
	return content, nil
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true。
func (s ConfigService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
