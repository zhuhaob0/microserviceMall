package service

import (
	"final-design/sk-admin/model"
	"log"

	"github.com/gohouse/gorose/v2"
)

type ProductService interface {
	CreateProduct(product *model.Product) error
	GetProductList() ([]gorose.Data, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(product *model.Product) error
}

type ProductServiceMiddleware func(ProductService) ProductService

type ProductServiceImpl struct{}

func (p ProductServiceImpl) CreateProduct(product *model.Product) error {
	productEntity := model.NewProductModel()
	err := productEntity.CreateProduct(product)
	if err != nil {
		log.Printf("ProductEntity.CreateProduct, err: %v", err)
		return err
	}
	return nil
}

func (p ProductServiceImpl) GetProductList() ([]gorose.Data, error) {
	productEntity := model.NewProductModel()
	productList, err := productEntity.GetProductList()
	if err != nil {
		log.Printf("productEntity.GetProductList, err: %v", err)
		return nil, err
	}
	return productList, nil
}

func (p ProductServiceImpl) UpdateProduct(product *model.Product) error {
	productEntity := model.NewProductModel()
	err := productEntity.UpdateProduct(product)
	if err != nil {
		log.Printf("ProductEntity.UpdateProduct, err: %v", err)
		return err
	}
	return nil
}

func (p ProductServiceImpl) DeleteProduct(product *model.Product) error {
	productEntity := model.NewProductModel()
	err := productEntity.DeleteProduct(product)
	if err != nil {
		log.Printf("ProductEntity.DeleteProduct, err: %v", err)
		return err
	}
	return nil
}
