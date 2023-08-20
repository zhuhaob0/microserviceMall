package plugins

import (
	"context"
	"errors"
	"final-design/pb"
	"final-design/pkg/client"
	"final-design/sk-app/model"
	"log"

	"github.com/go-kit/kit/endpoint"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
)

type RemoteOAuthService struct {
	oauthClient client.OAuthClient
}

func NewRemoteOAuthService() *RemoteOAuthService {
	oauthClient, _ := client.NewOAuthClient("oauth", nil, nil)
	return &RemoteOAuthService{
		oauthClient: oauthClient,
	}
}

func AuthToken() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(model.SecRequest)

			// log.Printf("req.token=%v\n", req.AccessToken)
			oauthService := NewRemoteOAuthService()
			resp, _ := oauthService.oauthClient.CheckToken(ctx, nil, &pb.CheckTokenRequest{
				Token: req.AccessToken,
			})
			if resp == nil {
				return nil, ErrTokenInvalid
			}
			if !resp.IsValidToken { // token无效，可能过期或者解析失败
				log.Printf("resp.Error = %v", resp.Err)
				return nil, errors.New(resp.Err)
			}
			log.Println("secKill的token鉴权成功")
			return next(ctx, request)
		}
	}
}
