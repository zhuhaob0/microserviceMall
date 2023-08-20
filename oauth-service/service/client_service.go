package service

import (
	"context"
	"errors"
	"final-design/oauth-service/model"
	"fmt"
)

var ErrClientMessage = errors.New("invalid client")

// Service Define a service interface
type ClientDetailsService interface {
	GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error)
}

type MysqlClientDetailsService struct{}

func NewMysqlClientDetailsService() ClientDetailsService {
	return &MysqlClientDetailsService{}
}

func (service *MysqlClientDetailsService) GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string) (*model.ClientDetails, error) {
	clientDetailsModel := model.NewClientDetailsModel()
	if clientDetails, err := clientDetailsModel.GetClientDetailsByClientId(clientId); err == nil {
		fmt.Println(clientDetails)
		if clientSecret == clientDetails.ClientSecret {
			return clientDetails, nil
		} else {
			return nil, ErrClientMessage
		}
	} else {
		fmt.Println("err:", err)
		return nil, err
	}
}
