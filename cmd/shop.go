package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/AlecAivazis/survey/v2"

	"github.com/siddmoitra/btech-minicash-2-poc/backend/domain"
	"github.com/siddmoitra/btech-minicash-2-poc/backend/handler"
	"github.com/siddmoitra/btech-minicash-2-poc/utils"
)

func questionSelectProducts(productNames []string) ([]domain.Product, error) {
	questionSelectProducts := []*survey.Question{
		{
			Name: "selectedProductNames",
			Prompt: &survey.MultiSelect{
				Message:  "Select product(s) to add to the shopping basket ?",
				Options:  productNames,
				PageSize: len(productNames),
			},
		},
	}

	// ask the question
	selectedProductNames := []string{}
	err := survey.Ask(questionSelectProducts, &selectedProductNames)

	if err != nil {
		return nil, err
	}

	selectedProducts := make([]domain.Product, len(selectedProductNames))
	// convert the product names to product IDs and return
	for i, selectedProduct := range selectedProductNames {
		product, err := handler.GetProductByName(selectedProduct)
		if err != nil {
			continue
		}
		selectedProducts[i] = product
	}
	return selectedProducts, nil
}

func questionCheckout() (bool, error) {
	shouldCheckout := false
	questionCheckout := []*survey.Question{
		{
			Name: "shouldCheckout",
			Prompt: &survey.Confirm{
				Message: "Do you want to checkout ?",
			},
		},
	}
	err := survey.Ask(questionCheckout, &shouldCheckout)
	if err != nil {
		return shouldCheckout, err
	}

	return shouldCheckout, nil
}

func questionDeliveryAddress() (string, error) {
	deliveryAddress := ""
	questionDeliveryAddress := []*survey.Question{
		{
			Name: "deliveryAddress",
			Prompt: &survey.Input{
				Message: "Enter the delivery address ?",
			},
		},
	}
	err := survey.Ask(questionDeliveryAddress, &deliveryAddress)
	if err != nil {
		return deliveryAddress, err
	}

	return deliveryAddress, nil
}

var payByMiniCash = "Pay by MiniCash 2.0"
var payByCard = "Pay by Card (not in scope of this demo)"

func questionPaymentMethod() (string, error) {
	paymentMethod := ""
	questionPaymentMethod := []*survey.Question{
		{
			Name: "selectedPaymentMethod",
			Prompt: &survey.Select{
				Message: "Select the payment method ?",
				Options: []string{payByCard, payByMiniCash},
			},
		},
	}
	err := survey.Ask(questionPaymentMethod, &paymentMethod)
	if err != nil {
		return paymentMethod, err
	}

	return paymentMethod, nil
}

func createSalesOrder(shippingAddress string, selectedProducts []domain.Product) *domain.SalesOrder {
	var salesOrder = new(domain.SalesOrder)

	salesOrder.ID = fmt.Sprintf("SO-%d", time.Now().Unix())
	salesOrder.ShippingAddress = shippingAddress
	salesOrder.LineItems = make([]domain.LineItem, len(selectedProducts))

	totalOrderPrice := 0.0
	for i, p := range selectedProducts {
		salesOrder.LineItems[i] = domain.LineItem{
			Product:  p,
			Quantity: 1,
		}
		totalOrderPrice += p.Price
	}
	salesOrder.TotalOrderPrice = totalOrderPrice

	return salesOrder
}

func ECommerceStartShopping() (*domain.SalesOrder, error) {

	utils.ShowUserInSystem("End-Consumer in e-Commerce shop (Magento)")

	// 1. List of products
	productNames := handler.ListProductNames()

	selectedProducts, err := questionSelectProducts(productNames)
	if err != nil {
		return nil, err
	}
	shouldCheckout, err := questionCheckout()
	if err != nil {
		return nil, err
	}
	if !shouldCheckout {
		log.Println("-> You chose not to checkout, will exit.")
		return nil, nil
	}

	deliveryAdress, err := questionDeliveryAddress()
	if err != nil {
		return nil, err
	}

	paymentMethod, err := questionPaymentMethod()
	if err != nil {
		return nil, err
	}

	if paymentMethod != payByMiniCash {
		log.Println("-> You chose pay by card which is not part of the demo, will exit.")
		return nil, nil
	}

	salesOrder := createSalesOrder(deliveryAdress, selectedProducts)

	return salesOrder, nil
}

func ECommerceHandlePaymentStatus(salesOrder *domain.SalesOrder, paymentStatus string) {
	fmt.Println("Payment Status Handled")
}
