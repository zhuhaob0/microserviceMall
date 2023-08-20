package service

import (
	"context"
	"errors"
	"final-design/oauth-service/model"
	"final-design/pb"
	"final-design/pkg/client"
)

var (
	ErrInvalidAuthentication = errors.New("invalid auth")
	ErrInvalidUserInfo       = errors.New("invalid user info")
)

// Service Define a service interface
type UserDetailsService interface {
	// Get UserDetails By Username
	GetUserDetailByUsername(ctx context.Context, username, password string, idType int) (*model.UserDetails, error)
}

// UserService implement Service interface
type RemoteUserService struct {
	userClient client.UserClient
}

func (service *RemoteUserService) GetUserDetailByUsername(ctx context.Context, username string, password string, idType int) (*model.UserDetails, error) {
	var (
		response *pb.UserResponse
		err      error
	)
	if idType == 0 {
		response, err = service.userClient.CheckUser(ctx, nil, &pb.UserRequest{
			Username: username,
			Password: password,
		})
	} else if idType == 1 {
		response, err = service.userClient.CheckAdminUser(ctx, nil, &pb.UserRequest{
			Username: username,
			Password: password,
		})
	}

	if err == nil {
		if response.UserId != 0 {
			return &model.UserDetails{
				UserId:   response.UserId,
				Username: username,
				Password: password,
			}, nil
		} else {
			return nil, ErrInvalidUserInfo
		}
	}
	return nil, err
}

func NewRemoteUserDetailService() *RemoteUserService {
	userClient, _ := client.NewUserClient("user", nil, nil)
	return &RemoteUserService{
		userClient: userClient,
	}
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
