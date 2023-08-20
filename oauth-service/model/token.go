package model

import "time"

type OAuth2Token struct {
	RefreshToken *OAuth2Token // 刷新令牌
	TokenType    string       // 令牌类型
	TokenValue   string       // 令牌
	ExpiresTime  *time.Time   // 过期时间
}

func (OAuth2Token *OAuth2Token) IsExpired() bool {
	return OAuth2Token.ExpiresTime != nil && OAuth2Token.ExpiresTime.Before(time.Now())
}

type OAuth2Details struct {
	Client *ClientDetails
	User   *UserDetails
}
