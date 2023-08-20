package model

import (
	"encoding/json"
	"final-design/pkg/mysql"
	"fmt"
	"log"
)

type ClientDetails struct {
	ClientId                    string   // client 的标识
	ClientSecret                string   // client 的密钥
	AccessTokenValiditySeconds  int      // 访问令牌有效时间：秒
	RefreshTokenValiditySeconds int      // 刷新令牌有效时间：秒
	RegisteredRedirectUri       string   // 重定向地址，授权码类型中使用
	AuthorizedGrantTypes        []string // 可以使用的授权类型
}

func (clientDetails *ClientDetails) IsMatch(clientId string, clientSecret string) bool {
	return clientId == clientDetails.ClientId && clientSecret == clientDetails.ClientSecret
}

type ClientDetailsModel struct{}

func NewClientDetailsModel() *ClientDetailsModel {
	return &ClientDetailsModel{}
}

func (c *ClientDetailsModel) getTableName() string {
	return "client_details"
}

func (c *ClientDetailsModel) GetClientDetailsByClientId(clientId string) (*ClientDetails, error) {
	conn := mysql.DB()
	// fmt.Println("GetClientDetailsByClientId:", clientId)
	result, err := conn.Table(c.getTableName()).Where(map[string]interface{}{"client_id": clientId}).First()
	fmt.Println("result:", result, "\nerr:", err)
	if err == nil {
		var authorizedGrantTypes []string
		_ = json.Unmarshal([]byte(result["authorized_grant_types"].(string)), &authorizedGrantTypes)

		return &ClientDetails{
			ClientId:                    result["client_id"].(string),
			ClientSecret:                result["client_secret"].(string),
			AccessTokenValiditySeconds:  int(result["access_token_validity_seconds"].(int64)),
			RefreshTokenValiditySeconds: int(result["refresh_token_validity_seconds"].(int64)),
			RegisteredRedirectUri:       result["registerd_redirect_uri"].(string),
			AuthorizedGrantTypes:        authorizedGrantTypes,
		}, nil
	} else {
		fmt.Println("查询client_id出错")
		return nil, err
	}
}

func (c *ClientDetailsModel) CreateClientDetails(clientDetails *ClientDetails) error {
	conn := mysql.DB()

	grantTypeString, _ := json.Marshal(clientDetails.AuthorizedGrantTypes)
	_, err := conn.Table(c.getTableName()).Data(map[string]interface{}{
		"client_id":                      clientDetails.ClientId,
		"client_secret":                  clientDetails.ClientSecret,
		"access_token_validity_seconds":  clientDetails.AccessTokenValiditySeconds,
		"refresh_token_validity_seconds": clientDetails.RefreshTokenValiditySeconds,
		"register_redirect_uri":          clientDetails.RegisteredRedirectUri,
		"authorized_grant_types":         grantTypeString,
	}).Insert()

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	return nil
}
