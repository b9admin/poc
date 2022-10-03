package handler

import (
	"fmt"

	"github.com/gocarina/gocsv"
	"github.com/siddmoitra/btech-minicash-2-poc/backend/domain"
	"github.com/siddmoitra/btech-minicash-2-poc/utils"
)

var products []*domain.Product

func ECommerceLoadData() {
	if products == nil {
		bytes := utils.LoadCSV("products.csv")

		if err := gocsv.UnmarshalString(bytes, &products); err != nil { // Load products from file
			panic(err)
		}
	}
}

func ListProducts() []*domain.Product {
	return products
}

func ListProductNames() []string {
	productNames := make([]string, len(products))
	for i, product := range products {
		productNames[i] = product.Name
	}
	return productNames
}

func GetProductByName(name string) (domain.Product, error) {
	for _, v := range products {
		if v.Name == name {
			return *v, nil
		}
	}

	return domain.Product{}, fmt.Errorf("no product by name: %s", name)
}
