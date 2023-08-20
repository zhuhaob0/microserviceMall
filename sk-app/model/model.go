package model

type SecRequest struct {
	ProductId     int             `json:"product_id"`     // 商品ID
	ProductName   string          `json:"product_name"`   // 商品名
	ActivityPrice int             `json:"activity_price"` // 活动价
	Source        string          `json:"source"`
	AuthCode      string          `json:"auth_code"`
	SecTime       int64           `json:"sec_time"` // 秒杀时间
	Nance         string          `json:"nance"`
	UserId        int             `json:"user_id"`
	Username      string          `json:"username"`
	UserAuthSign  string          `json:"user_auth_sign"` // 用户授权签名
	AccessTime    int64           `json:"access_time"`    // 访问时间
	AccessToken   string          `json:"access_token"`   // 访问令牌
	ClientAddr    string          `json:"client_addr"`
	ClientRefence string          `json:"client_refence"`
	CloseNotify   <-chan bool     `json:"-"`
	ResultChan    chan *SecResult `json:"-"`
}

type SecResult struct {
	ProductId int    `json:"product_id"` // 商品ID
	UserId    int    `json:"user_id"`    // 用户ID
	Token     string `json:"token"`      // Token
	TokenTime int64  `json:"token_time"` // Token生成时间
	Code      int    `json:"code"`       // 状态码
}

type Order struct {
	OrderId       int    `json:"order_id"`       // 订单ID
	ProductId     int    `json:"product_id"`     // 购买的商品ID
	ProductName   string `json:"product_name"`   // 购买的商品名
	OrderTime     int64  `json:"order_time"`     // 下单时间
	Buyer         string `json:"buyer"`          // 买家
	ActivityPrice int    `json:"activity_price"` // 活动价
}
