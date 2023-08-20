package model

import (
	"errors"
	"final-design/pkg/mysql"
	"fmt"
	"log"

	"github.com/gohouse/gorose/v2"
	"github.com/unknwon/com"
)

type Product struct {
	ProductId   int    `json:"product_id"`   // 商品Id
	ProductName string `json:"product_name"` // 商品名称
	Price       int    `json:"price"`        // 价格
	ImgUri      string `json:"img_uri"`      // 图片存储位置
	Description string `json:"description"`  // 描述
}

type ProductModel struct{}

func NewProductModel() *ProductModel {
	return &ProductModel{}
}

func (p *ProductModel) getTableName() string {
	return "product"
}

func (p *ProductModel) GetProductList() ([]gorose.Data, error) {
	conn := mysql.DB()
	list, err := conn.Table(p.getTableName()).Get()
	if err != nil {
		log.Printf("GetProductList, Error: %v", err)
		return nil, err
	}
	return list, err
}

func (p *ProductModel) CreateProduct(product *Product) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"product_name": product.ProductName,
		"price":        product.Price,
		"img_uri":      product.ImgUri,
		"description":  product.Description,
	}).Insert()
	if err != nil {
		log.Printf("CreateProduct, Error: %v", err)
		return err
	}
	return nil
}

func (p *ProductModel) UpdateProduct(product *Product) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"price":       product.Price,
		"img_uri":     product.ImgUri,
		"description": product.Description,
	}).Where("product_name", product.ProductName).Update()
	if err != nil {
		log.Printf("UpdateProduct, Error: %v\n", err)
		return err
	}
	return nil
}

func (p *ProductModel) DeleteProduct(product *Product) error {
	conn := mysql.DB()
	_, err := conn.Table(p.getTableName()).Where("product_name", product.ProductName).Delete()
	if err != nil {
		log.Printf("DeleteProduct, Error: %v\n", err)
		return err
	}
	return nil
}

func (p *ProductModel) GetProductIdByName(productName string) (int, error) {
	conn := mysql.DB()
	product, err := conn.Table(p.getTableName()).Where("product_name", productName).First()
	productId, _ := com.StrTo(fmt.Sprint(product["product_id"])).Int()
	fmt.Println("product=", product)
	if productId == 0 || err != nil {
		log.Println("GetProductIdByName, 内部查询出错或者商品不存在")
		return 0, errors.New("商品不存在")
	}
	return productId, nil
}
