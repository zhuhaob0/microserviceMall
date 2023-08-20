package service

import (
	"context"
	"final-design/user-service/model"
	"log"
)

// Service define a service interface
type Service interface {
	Create(ctx context.Context, username, password string, userId int64, age int) error
	Check(ctx context.Context, username, password string) (int64, error)

	CreateAdmin(ctx context.Context, username, password string, userId int64, age int) error
	CheckAdmin(ctx context.Context, username string, password string) (int64, error)

	// HealthCheck check service health status
	HealthCheck() bool
}

// UserService implement Service interface
type UserService struct{}

func (s UserService) Create(ctx context.Context, username, password string, userId int64, age int) error {
	userEntity := model.NewUserModel()
	err := userEntity.CreateUser(&model.User{
		Username: username,
		Password: password,
		UserId:   userId,
		Age:      age,
	})
	return err
}

func (s UserService) Check(ctx context.Context, username string, password string) (int64, error) {
	userEntity := model.NewUserModel()
	res, err := userEntity.CheckUser(username, password)
	if err != nil {
		log.Printf("UserEntity.CheckUser, err: %v", err)
		return 0, err
	}
	return res.UserId, nil
}

func (s UserService) CreateAdmin(ctx context.Context, username, password string, userId int64, age int) error {
	userEntity := model.NewAdminUserModel()
	err := userEntity.CreateUser(&model.AdminUser{
		Username: username,
		Password: password,
		UserId:   userId,
		Age:      age,
	})
	return err
}

func (s UserService) CheckAdmin(ctx context.Context, username string, password string) (int64, error) {
	userEntity := model.NewAdminUserModel()
	res, err := userEntity.CheckUser(username, password) 
	if err != nil {
		log.Printf("UserEntity.CheckUser, err: %v", err)
		return 0, err
	}
	return res.UserId, nil
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s UserService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
