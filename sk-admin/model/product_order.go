package model

import (
	"final-design/pkg/mysql"
	"log"

	"github.com/gohouse/gorose/v2"
)

type Order struct {
	OrderId       int    `json:"order_id"`       // 订单ID
	ProductId     int    `json:"product_id"`     // 购买的商品ID
	ProductName   string `json:"product_name"`   // 购买的商品名
	OrderTime     int64  `json:"order_time"`     // 下单时间
	Buyer         string `json:"buyer"`          // 买家
	ActivityPrice int    `json:"activity_price"` // 活动价
}

type OrderModel struct{}

func NewOrderModel() *OrderModel {
	return &OrderModel{}
}

func (o *OrderModel) getTableName() string {
	return "product_order"
}

func (p *OrderModel) GetOrderList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Order("product_id asc").Get()
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	return list, nil
}

func (p *OrderModel) CreateOrder(order *Order) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"product_id":     order.ProductId,
		"product_name":   order.ProductName,
		"order_time":     order.OrderTime,
		"buyer":          order.Buyer,
		"activity_price": order.ActivityPrice,
	}).Insert()
	if err != nil {
		log.Printf("CreateOrder, Error: %v", err)
		return err
	}
	return nil
}

func (p *OrderModel) GetBuyerOrder(buyer string) ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Fields("product_name, order_time, activity_price").Where("buyer", "=", buyer).Get()
	if err != nil {
		log.Printf("GetBuyerOrder, Error: %v", err)
		return nil, err
	}
	return list, nil
}
