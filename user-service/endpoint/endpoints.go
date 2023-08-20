package endpoint

import (
	"context"
	"errors"
	"final-design/user-service/service"

	"github.com/go-kit/kit/endpoint"
)

var (
	ErrInvalidRequestType = errors.New("invalid username, password")
)

type UserEndpoints struct {
	CreateUserEndpoint endpoint.Endpoint
	UserEndpoint       endpoint.Endpoint

	CreateAdminUserEndpoint endpoint.Endpoint
	AdminUserEndpoint       endpoint.Endpoint

	HealthCheckEndpoint endpoint.Endpoint
}

// UserEndpoint define endpoint
func (u *UserEndpoints) Check(ctx context.Context, username, password string) (int64, error) {
	// reflect.TypeOf(UserEndpoints{})
	resp, err := u.UserEndpoint(ctx, UserRequest{
		Username: username,
		Password: password,
	})
	response := resp.(UserResponse)
	if err != nil {
		err = errors.New("bad request")
	}
	return response.UserId, err
}

// UserEndpoint define endpoint
func (u *UserEndpoints) AdminCheck(ctx context.Context, username, password string) (int64, error) {
	// reflect.TypeOf(UserEndpoints{})
	resp, err := u.AdminUserEndpoint(ctx, UserRequest{
		Username: username,
		Password: password,
	})
	response := resp.(UserResponse)
	if err != nil {
		err = errors.New("bad request")
	}
	return response.UserId, err
}

func (u *UserEndpoints) Create(ctx context.Context, username, password string, userId int64, age int) error {
	// reflect.TypeOf(UserEndpoints{})
	resp, _ := u.CreateUserEndpoint(ctx, CreateUserRequest{
		Username: username,
		Password: password,
		UserId:   userId,
		Age:      age,
	})
	response := resp.(CreateUserResponse)
	return response.Error
}

func (u *UserEndpoints) HealthCheck() bool {
	return false
}

// UserRequest define request struct
type UserRequest struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

// UserResponse define response struct
type UserResponse struct {
	Result bool  `json:"result"`
	UserId int64 `json:"user_id"`
	Error  error `json:"error"`
}

// 创建检查User的endpoint
func MakeUserEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserRequest)
		var (
			username, password string
			userId             int64
			calError           error
		)
		username = req.Username
		password = req.Password
		userId, calError = svc.Check(ctx, username, password)
		if calError != nil {
			return UserResponse{Result: false, Error: calError}, nil
		}
		return UserResponse{Result: true, UserId: userId, Error: calError}, nil
	}
}

// 创建检查AdminUser的endpoint
func MakeAdminUserEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserRequest)
		var (
			username, password string
			userId             int64
			calError           error
		)
		username = req.Username
		password = req.Password
		userId, calError = svc.CheckAdmin(ctx, username, password)
		if calError != nil {
			return UserResponse{Result: false, Error: calError}, nil
		}
		return UserResponse{Result: true, UserId: userId, Error: calError}, nil
	}
}

type CreateUserRequest struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
	UserId   int64  `json:"user_id"`
	Age      int    `json:"age"`
}

type CreateUserResponse struct {
	Result bool  `json:"result"`
	Error  error `json:"error"`
}

// 创建注册User的endpoint
func MakeCreateUserEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateUserRequest)
		var (
			username, password string
			userId             int64
			calError           error
			age                int
		)
		username, password, userId, age = req.Username, req.Password, req.UserId, req.Age
		calError = svc.Create(ctx, username, password, userId, age)
		return CreateUserResponse{Result: calError == nil, Error: calError}, nil
	}
}

// 创建注册AdminUser的endpoint
func MakeCreateAdminUserEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateUserRequest)
		var (
			username, password string
			userId             int64
			calError           error
			age                int
		)
		username, password, userId, age = req.Username, req.Password, req.UserId, req.Age
		calError = svc.CreateAdmin(ctx, username, password, userId, age)
		return CreateUserResponse{Result: calError == nil, Error: calError}, nil
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
