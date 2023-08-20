package service

import (
	"encoding/json"
	"final-design/sk-admin/model"
	"time"

	"log"

	"github.com/gohouse/gorose/utils"
	"github.com/gohouse/gorose/v2"
)

type OrderService interface {
	GetOrderList() (map[string]interface{}, error)
	GetBuyerOrder(buyer string) ([]gorose.Data, error)
}

type OrderServiceMiddleware func(OrderService) OrderService

type OrderServiceImpl struct{}

func (o *OrderServiceImpl) GetOrderList() (map[string]interface{}, error) {
	orderEntity := model.NewOrderModel()
	orderList, err := orderEntity.GetOrderList()
	if err != nil {
		log.Printf("orderEntity.GetOrderList, err: %v", err)
		return nil, err
	}
	result := make(map[string][4]int, 10)
	nowTime := time.Now().Unix()
	hour := int64(3600) // 一小时3600秒
	for _, v := range orderList {
		var tmp model.Order
		jsonStr, _ := utils.JsonEncode(v)
		json.Unmarshal([]byte(jsonStr), &tmp)
		val := result[tmp.ProductName]
		if nowTime-tmp.OrderTime < hour { // 一小时内的订单
			val[0]++
		}
		if nowTime-tmp.OrderTime < hour*24 { // 一天内的订单
			val[1]++
		}
		if nowTime-tmp.OrderTime < hour*24*7 { // 一周内的订单
			val[2]++
		}
		val[3] = val[1] * tmp.ActivityPrice
		result[tmp.ProductName] = val
	}
	ret := make(map[string]interface{}, 10)
	for k, v := range result {
		ret[k] = v
	}
	return ret, nil
}

func (o *OrderServiceImpl) GetBuyerOrder(buyer string) ([]gorose.Data, error) {
	orderEntity := model.NewOrderModel()
	buyerOrder, err := orderEntity.GetBuyerOrder(buyer)
	if err != nil {
		log.Printf("orderEntity.GetBuyerOrder, err: %v", err)
		return nil, err
	}
	return buyerOrder, nil
}
