package model

import (
	"final-design/pkg/mysql"
	"fmt"
	"log"

	"github.com/gohouse/gorose/v2"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type Activity struct {
	ActivityId    int    `json:"activity_id"`    // 活动Id
	ActivityName  string `json:"activity_name"`  // 活动名称
	ProductId     int    `json:"product_id"`     // 商品Id
	ProductName   string `json:"product_name"`   // 商品名
	StartTime     int64  `json:"start_time"`     // 开始时间
	EndTime       int64  `json:"end_time"`       // 结束时间
	Total         int    `json:"total"`          // 商品总数
	LeftNum       int    `json:"left_num"`       // 剩余数量
	Status        int    `json:"status"`         // 状态
	ActivityPrice int    `json:"activity_price"` // 活动价

	StartTimeStr     string  `json:"start_time_str"`
	EndTimeStr       string  `json:"end_time_str"`
	StatusStr        string  `json:"status_str"`
	MaxSoldPerSecond int     `json:"max_sold_per_second"`
	MaxBuyPerPerson  int     `json:"max_buy_per_person"`
	BuyRate          float64 `json:"buy_rate"`
}

type SecProductInfoConf struct {
	ActivityName     string  `json:"activity_name"`       // 活动名
	ProductId        int     `json:"product_id"`          // 商品Id
	ProductName      string  `json:"product_name"`        // 商品名
	StartTime        int64   `json:"start_time"`          // 开始时间
	EndTime          int64   `json:"end_time"`            // 结束时间
	ActivityPrice    int     `json:"activity_price"`      // 活动价
	Status           int     `json:"status"`              // 状态
	Total            int     `json:"total"`               // 商品总数
	LeftNum          int     `json:"left_num"`            // 剩余商品数
	MaxBuyPerPerson  int     `json:"max_buy_per_person"`  // 单人购买限制
	MaxSoldPerSecond int     `json:"max_sold_per_second"` // 每秒最多能卖多少
	BuyRate          float64 `json:"buy_rate"`            // 买中几率
}

type ActivityModel struct{}

func NewActivityModel() *ActivityModel {
	return &ActivityModel{}
}

func (p *ActivityModel) getTableName() string {
	return "activity"
}

func (p *ActivityModel) GetActivityList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Order("end_time desc").Get()
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	return list, nil
}

func (p *ActivityModel) CreateActivity(activity *Activity) error {
	conn := mysql.DB()
	productModel := NewProductModel()
	productId, err := productModel.GetProductIdByName(activity.ProductName)
	if err != nil { // 商品还未出现过，就创建它
		productModel.CreateProduct(&Product{
			ProductName: activity.ProductName,
			Price:       activity.ActivityPrice,
			ImgUri:      "default.jpg",        // 默认图片是default.jpg
			Description: activity.ProductName, // 商品描述默认为productName
		})
		productId, _ = productModel.GetProductIdByName(activity.ProductName)
	}

	activity.ProductId = productId
	_, err = conn.Table(p.getTableName()).Data(map[string]interface{}{
		"activity_name": activity.ActivityName,
		"product_id":    activity.ProductId,
		"product_name":  activity.ProductName,
		"start_time":    activity.StartTime,
		"end_time":      activity.EndTime,
		"total":         activity.Total,
		"left_num":      activity.LeftNum,
		"status":        activity.Status,

		"max_sold_per_second": activity.MaxSoldPerSecond,
		"max_buy_per_person":  activity.MaxBuyPerPerson,
		"activity_price":      activity.ActivityPrice,
		"buy_rate":            1,
	}).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (p *ActivityModel) UpdateActivity(activity *Activity) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"start_time":          activity.StartTime,
		"end_time":            activity.EndTime,
		"total":               activity.Total,
		"left_num":            activity.LeftNum,
		"status":              activity.Status,
		"max_sold_per_second": activity.MaxSoldPerSecond,
		"max_buy_per_person":  activity.MaxBuyPerPerson,
		"activity_price":      activity.ActivityPrice,
	}).Where("activity_name", activity.ActivityName).Update()
	if err != nil {
		fmt.Println("activity 更新失败")
		return err
	}
	return nil
}

func (p *ActivityModel) DeleteActivity(activity *Activity) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Where("activity_name", activity.ActivityName).Delete()
	if err != nil {
		return err
	}
	return nil
}
